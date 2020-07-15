from __future__ import unicode_literals

import os

import typing
from six import StringIO
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
    GoImportHelper)


class GoClientBackend(CodeBackend):


    def __init__(self, target_folder_path, args):
        # type: (str, typing.Optional[typing.Sequence[str]]) -> None
        super(GoClientBackend, self).__init__(target_folder_path, args)
        self.import_helper = GoImportHelper()

    def generate(self, api):
        for namespace in api.namespaces.values():
            if len(namespace.routes) > 0:
                self._generate_client(namespace)

    def _generate_client(self, namespace):
        file_name = os.path.join(self.target_folder_path, namespace.name,
                                 'client.go')
        with self.output_to_relative_path(file_name):
            self.import_helper.reset()
            self.emit_raw(HEADER)
            self.emit()
            self.emit('package %s' % namespace.name)
            self.emit()
            output_buffer = StringIO()
            with self.capture_emitted_output(output_buffer):

                self.emit('// Client interface describes all routes in this namespace')
                with self.block('type Client interface'):
                    for route in namespace.routes:
                        generate_doc(self, route)
                        self.emit(self._generate_route_signature(namespace, route))
                self.emit()
                dropbox_imp = self.import_helper.id_for_package("github.com/dropbox/dropbox-sdk-go-unofficial/dropbox")

                self.emit('type apiImpl %s.Context' %dropbox_imp)
                for route in namespace.routes:
                    self._generate_route(namespace, route)
                self.emit('// New returns a Client implementation for this namespace')
                with self.block('func New(c %s.Config) Client' %dropbox_imp):
                    self.emit('ctx := apiImpl(%s.NewContext(c))' %dropbox_imp)
                    self.emit('return &ctx')

            self.import_helper.emit_import_statements(self)

            self._append_output(output_buffer.getvalue())


    def _generate_route_signature(self, namespace, route):
        req = fmt_type(self.import_helper, route.arg_data_type, namespace)
        res = fmt_type(self.import_helper, route.result_data_type, namespace, use_interface=True)
        fn = fmt_var(route.name)
        if route.version != 1:
            fn += 'V%d' % route.version
        style = route.attrs.get('style', 'rpc')

        arg = '' if is_void_type(route.arg_data_type) else 'arg {req}'
        ret = '(err error)' if is_void_type(route.result_data_type) else \
            '(res {res}, err error)'
        signature = '{fn}(' + arg + ') ' + ret
        if style == 'download':
            io_imp = self.import_helper.id_for_package(
                "io")

            signature = '{fn}(' + arg + \
                ') (res {res}, content '+io_imp+'.ReadCloser, err error)'
        elif style == 'upload':
            io_imp = self.import_helper.id_for_package(
                "io")
            signature = '{fn}(' + arg + ', content '+io_imp+'.Reader) ' + ret
            if is_void_type(route.arg_data_type):
                signature = '{fn}(content io.Reader) ' + ret
        return signature.format(fn=fn, req=req, res=res)


    def _generate_route(self, namespace, route):
        out = self.emit
        fn = fmt_var(route.name)
        if route.version != 1:
            fn += 'V%d' % route.version
        err = fmt_type(self.import_helper, route.error_data_type, namespace)
        out('//%sAPIError is an error-wrapper for the %s route' %
            (fn, route.name))
        with self.block('type {fn}APIError struct'.format(fn=fn)):
            dropbox_imp = self.import_helper.id_for_package(
                "github.com/dropbox/dropbox-sdk-go-unofficial/dropbox")

            out('%s.APIError'%dropbox_imp)
            out('EndpointError {err} `json:"error"`'.format(err=err))
        out()

        signature = 'func (dbx *apiImpl) ' + self._generate_route_signature(
            namespace, route)
        with self.block(signature):
            if route.deprecated is not None:
                log_imp = self.import_helper.id_for_package(
                    "log")
                out('%s.Printf("WARNING: API `%s` is deprecated")' % (log_imp, fn))
                if route.deprecated.by is not None:
                    out('%s.Printf("Use API `%s` instead")' % (log_imp, fmt_var(route.deprecated.by.name)))
                out()

            out('cli := dbx.Client')
            out()

            self._generate_request(namespace, route)
            self._generate_post()
            self._generate_response(route)
            http_imp = self.import_helper.id_for_package("net/http")

            ok_check = 'if resp.StatusCode == %s.StatusOK' %http_imp
            if fn == "Download":
                ok_check += ' || resp.StatusCode == %s.StatusPartialContent' %http_imp
            with self.block(ok_check):
                self._generate_result(route)
            self._generate_error_handling(namespace, route)

        out()

    def _generate_request(self, namespace, route):
        out = self.emit
        auth = route.attrs.get('auth', '')
        host = route.attrs.get('host', 'api')
        style = route.attrs.get('style', 'rpc')

        body = 'nil'
        if not is_void_type(route.arg_data_type):
            out('dbx.Config.LogDebug("arg: %v", arg)')
            json_imp = self.import_helper.id_for_package("encoding/json")

            out('b, err := %s.Marshal(arg)'%json_imp)
            with self.block('if err != nil'):
                out('return')
            out()
            if host != 'content':
                bytes_imp = self.import_helper.id_for_package("bytes")
                body = '%s.NewReader(b)'%bytes_imp
        if style == 'upload':
            body = 'content'

        headers = {}
        if not is_void_type(route.arg_data_type):
            if host == 'content' or style in ['upload', 'download']:
                dropbox_imp = self.import_helper.id_for_package("github.com/dropbox/dropbox-sdk-go-unofficial/dropbox")

                headers["Dropbox-API-Arg"] = "%s.HTTPHeaderSafeJSON(b)" % dropbox_imp
            else:
                headers["Content-Type"] = '"application/json"'
        if style == 'upload':
            headers["Content-Type"] = '"application/octet-stream"'

        out("headers := map[string]string{")
        for k, v in sorted(headers.items()):
            out('\t"{}": {},'.format(k, v))
        out("}")
        if fmt_var(route.name) == "Download":
            out("for k, v := range arg.ExtraHeaders { headers[k] = v }")
        if auth != 'noauth' and auth != 'team':
            with self.block('if dbx.Config.AsMemberID != ""'):
                out('headers["Dropbox-API-Select-User"] = dbx.Config.AsMemberID')
        out()

        fn = route.name
        if route.version != 1:
            fn += '_v%d' % route.version
        authed = 'false' if auth == 'noauth' else 'true'
        dropbox_imp = self.import_helper.id_for_package(
            "github.com/dropbox/dropbox-sdk-go-unofficial/dropbox")

        out('req, err := (*{}.Context)(dbx).NewRequest("{}", "{}", {}, "{}", "{}", headers, {})'.format(
            dropbox_imp, host, style, authed, namespace.name, fn, body))
        with self.block('if err != nil'):
            out('return')

        out('dbx.Config.LogInfo("req: %v", req)')

        out()

    def _generate_post(self):
        out = self.emit

        out('resp, err := cli.Do(req)')

        with self.block('if err != nil'):
            out('return')
        out()

        out('dbx.Config.LogInfo("resp: %v", resp)')

    def _generate_response(self, route):
        out = self.emit
        style = route.attrs.get('style', 'rpc')
        if style == 'download':
            out('body := []byte(resp.Header.Get("Dropbox-API-Result"))')
            out('content = resp.Body')
        else:
            out('defer resp.Body.Close()')
            ioutil_imp = self.import_helper.id_for_package("io/ioutil")

            with self.block('body, err := %s.ReadAll(resp.Body);'
                            'if err != nil'%ioutil_imp):
                out('return')
            out()

        out('dbx.Config.LogDebug("body: %s", body)')

    def _generate_error_handling(self, namespace, route):
        out = self.emit
        style = route.attrs.get('style', 'rpc')
        http_imp = self.import_helper.id_for_package("net/http")

        with self.block('if resp.StatusCode == %s.StatusConflict'%http_imp):
            # If style was download, body was assigned to a header.
            # Need to re-read the response body to parse the error
            if style == 'download':
                out('defer resp.Body.Close()')
                ioutil_imp = self.import_helper.id_for_package("io/ioutil")

                with self.block('body, err = %s.ReadAll(resp.Body);'
                                'if err != nil'%ioutil_imp):
                    out('return')
            out('var apiError %sAPIError' % fmt_var(route.name))
            json_imp = self.import_helper.id_for_package("encoding/json")

            with self.block('err = %s.Unmarshal(body, &apiError);' 
                            'if err != nil'% json_imp):
                out('return')
            out('err = apiError')
            out('return')

        auth_ns = ""
        if namespace.name != "auth":
            auth_ns = "%s."%self.import_helper.id_for_package(
                "github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/auth")

        with self.block('err = %sHandleCommonAuthErrors(dbx.Config, resp, body);'
                        'if err != nil' % auth_ns):
            out('return')
        dropbox_imp = self.import_helper.id_for_package(
            "github.com/dropbox/dropbox-sdk-go-unofficial/dropbox")

        out('err = %s.HandleCommonAPIErrors(dbx.Config, resp, body)' % dropbox_imp)
        out('return')

    def _generate_result(self, route):
        out = self.emit
        if is_struct_type(route.result_data_type) and \
                route.result_data_type.has_enumerated_subtypes():
            out('var tmp %sUnion' % fmt_var(route.result_data_type.name, export=False))
            json_imp = self.import_helper.id_for_package("encoding/json")

            with self.block('err = %s.Unmarshal(body, &tmp);'
                            'if err != nil'%json_imp):
                out('return')
            with self.block('switch tmp.Tag'):
                for t in route.result_data_type.get_enumerated_subtypes():
                    with self.block('case "%s":' % t.name, delim=(None, None)):
                        self.emit('res = tmp.%s' % fmt_var(t.name))
        elif not is_void_type(route.result_data_type):
            json_imp = self.import_helper.id_for_package("encoding/json")

            with self.block('err = %s.Unmarshal(body, &res);'
                            'if err != nil' % json_imp):
                out('return')
            out()
        out('return')
