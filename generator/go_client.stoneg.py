import os

from stone.backend import CodeBackend
from stone.ir import (
    is_void_type,
    is_struct_type
)

from go_helpers import (
    HEADER,
    fmt_type,
    fmt_var,
    generate_doc,
)


class GoClientBackend(CodeBackend):
    def generate(self, api):
        for namespace in api.namespaces.values():
            if len(namespace.routes) > 0:
                self._generate_client(namespace)

    def _generate_client(self, namespace):
        file_name = os.path.join(self.target_folder_path, namespace.name,
                                 'client.go')
        with self.output_to_relative_path(file_name):
            self.emit_raw(HEADER)
            self.emit()
            self.emit('package %s' % namespace.name)
            self.emit()

            self.emit('// Client interface describes all routes in this namespace')
            with self.block('type Client interface'):
                for route in namespace.routes:
                    generate_doc(self, route)
                    self.emit(self._generate_route_signature(namespace, route))
            self.emit()

            self.emit('// ContextClient interface describes all routes in this namespace with context support')
            with self.block('type ContextClient interface'):
                self.emit('Client')
                for route in namespace.routes:
                    generate_doc(self, route, name_override=self._generate_route_name(route, context=True))
                    self.emit(self._generate_route_signature(namespace, route, context=True))
            self.emit()

            self.emit('type apiImpl dropbox.Context')
            for route in namespace.routes:
                self._generate_route(namespace, route)
            self.emit('// NewContext returns a ContextClient implementation for this namespace')
            with self.block('func NewContext(c dropbox.Config) ContextClient'):
                self.emit('ctx := apiImpl(dropbox.NewContext(c))')
                self.emit('return &ctx')
            self.emit()
            self.emit('// New returns a Client implementation for this namespace')
            with self.block('func New(c dropbox.Config) Client'):
                self.emit('return NewContext(c)')

    def _generate_route_signature(self, namespace, route, context=False):
        req = fmt_type(route.arg_data_type, namespace)
        res = fmt_type(route.result_data_type, namespace, use_interface=True)
        fn = self._generate_route_name(route, context=context)
        style = route.attrs.get('style', 'rpc')

        args = []
        if context:
            args.append('ctx context.Context')
        if not is_void_type(route.arg_data_type):
            args.append('arg %s' % req)
        if style == 'upload':
            args.append('content io.Reader')

        rets = []
        if not is_void_type(route.result_data_type):
            rets.append('res %s' % res)
        if style == 'download':
            rets.append('content io.ReadCloser')
        rets.append('err error')

        return '{fn}({args}) ({rets})'.format(
            fn=fn, args=', '.join(args), rets=', '.join(rets))

    def _generate_route_name(self, route, context=False):
        fn = fmt_var(route.name)
        if route.version != 1:
            fn += 'V%d' % route.version

        if context:
            fn += 'Context'
        return fn

    def _generate_route(self, namespace, route):
        out = self.emit

        route_name = route.name
        if route.version != 1:
            route_name += '_v%d' % route.version

        route_style = route.attrs.get('style', 'rpc')

        fn = fmt_var(route.name)
        if route.version != 1:
            fn += 'V%d' % route.version

        err = fmt_type(route.error_data_type, namespace)
        out('//%sAPIError is an error-wrapper for the %s route' %
            (fn, route_name))
        with self.block('type {fn}APIError struct'.format(fn=fn)):
            out('dropbox.APIError')
            out('EndpointError {err} `json:"error"`'.format(err=err))
        out()

        # Generate the Context variant (does the real work)
        generate_doc(self, route, name_override=self._generate_route_name(route, context=True))
        signature_context = 'func (dbx *apiImpl) ' + self._generate_route_signature(
            namespace, route, context=True)
        with self.block(signature_context):
            if route.deprecated is not None:
                out('log.Printf("WARNING: API `%s` is deprecated")' % fn)
                if route.deprecated.by is not None:
                    replacement_fn = fmt_var(route.deprecated.by.name)
                    if route.deprecated.by.version != 1:
                        replacement_fn += "V%d" % route.deprecated.by.version
                    out('log.Printf("Use API `%s` instead")' % replacement_fn)
                out()

            args = {
                "Host": route.attrs.get('host', 'api'),
                "Namespace": namespace.name,
                "Route": route_name,
                "Auth": route.attrs.get('auth', ''),
                "Style": route_style,
            }

            with self.block('req := dropbox.Request'):
                for k, v in args.items():
                    out(k + ':"' + v + '",')

                out("Arg: {arg},".format(arg="arg" if not is_void_type(route.arg_data_type) else "nil"))
                out("ExtraHeaders: {headers},".format(
                    headers="arg.ExtraHeaders" if fmt_var(route.name) == "Download" else "nil"))
            out()

            out("var resp []byte")
            out("var respBody io.ReadCloser")
            out("resp, respBody, err = (*dropbox.Context)(dbx).ExecuteContext(ctx, req, {body})".format(
                body="content" if route_style == 'upload' else "nil"))
            with self.block("if err != nil"):
                out("var appErr {fn}APIError".format(fn=fn))
                out("err = {auth}ParseError(err, &appErr)".format(
                    auth="auth." if namespace.name != "auth" else ""))
                with self.block("if errors.Is(err, &appErr)"):
                    out("err = appErr")
                out("return")
            out()

            if is_struct_type(route.result_data_type) and route.result_data_type.has_enumerated_subtypes():
                out('var tmp %sUnion' % fmt_var(route.result_data_type.name, export=False))
                with self.block('err = json.Unmarshal(resp, &tmp);'
                                'if err != nil'):
                    out('return')
                with self.block('switch tmp.Tag'):
                    for t in route.result_data_type.get_enumerated_subtypes():
                        with self.block('case "%s":' % t.name, delim=(None, None)):
                            self.emit('res = tmp.%s' % fmt_var(t.name))
            elif not is_void_type(route.result_data_type):
                with self.block('err = json.Unmarshal(resp, &res);'
                                'if err != nil'):
                    out('return')
                out()
            else:
                out("_ = resp")

            if route_style == "download":
                out("content = respBody")
            else:
                out("_ = respBody")
            out('return')
        out()

        # Generate the non-context variant (delegates to Context variant)
        signature = 'func (dbx *apiImpl) ' + self._generate_route_signature(
            namespace, route, context=False)
        with self.block(signature):
            call_args = ['context.Background()']
            if not is_void_type(route.arg_data_type):
                call_args.append('arg')
            if route_style == 'upload':
                call_args.append('content')
            out('return dbx.%sContext(%s)' % (fn, ', '.join(call_args)))
        out()
