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

            self.emit('type apiImpl dropbox.Context')
            for route in namespace.routes:
                self._generate_route(namespace, route)
            self.emit('// New returns a Client implementation for this namespace')
            with self.block('func New(c dropbox.Config) Client'):
                self.emit('ctx := apiImpl(dropbox.NewContext(c))')
                self.emit('return &ctx')

    def _generate_route_signature(self, namespace, route):
        req = fmt_type(route.arg_data_type, namespace)
        res = fmt_type(route.result_data_type, namespace, use_interface=True)
        fn = fmt_var(route.name)
        if route.version != 1:
            fn += 'V%d' % route.version
        style = route.attrs.get('style', 'rpc')

        arg = '' if is_void_type(route.arg_data_type) else 'arg {req}'
        ret = '(err error)' if is_void_type(route.result_data_type) else \
            '(res {res}, err error)'
        signature = '{fn}(' + arg + ') ' + ret
        if style == 'download':
            signature = '{fn}(' + arg + \
                ') (res {res}, content io.ReadCloser, err error)'
        elif style == 'upload':
            signature = '{fn}(' + arg + ', content io.Reader) ' + ret
            if is_void_type(route.arg_data_type):
                signature = '{fn}(content io.Reader) ' + ret
        return signature.format(fn=fn, req=req, res=res)


    def _generate_route(self, namespace, route):
        out = self.emit

        route_name = route.name
        if route.version != 1:
            route_name += '_v%d' % route.version

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

        signature = 'func (dbx *apiImpl) ' + self._generate_route_signature(
            namespace, route)
        with self.block(signature):
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
                "Style": route.attrs.get('style', 'rpc'),
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
            out("resp, respBody, err = (*dropbox.Context)(dbx).Execute(req, {body})".format(
                body="content" if route.attrs.get('style', '') == 'upload' else "nil"))
            with self.block("if err != nil"):
                out("var appErr {fn}APIError".format(fn=fn))
                out("err = {auth}ParseError(err, &appErr)".format(
                    auth="auth." if namespace.name != "auth" else ""))
                with self.block("if err == &appErr"):
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

            if route.attrs.get('style', 'rpc') == "download":
                out("content = respBody")
            else:
                out("_ = respBody")
            out('return')
        out()
