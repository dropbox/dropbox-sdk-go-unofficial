from __future__ import unicode_literals

from stone.backend import CodeBackend
from stone.ir import (ApiNamespace, ApiRoute)
from stone.ir import (
    Boolean,
    Float32,
    Float64,
    Int32,
    Int64,
    String,
    Timestamp,
    UInt32,
    UInt64,
    unwrap_nullable,
    is_composite_type,
    is_list_type,
    is_struct_type,
    Void,
)
from stone.backends import helpers
from typing import Text, Dict

HEADER = """\
// Copyright (c) Dropbox, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
"""

_reserved_keywords = {
    'break', 'default', 'func', 'interface', 'select',
    'case', 'defer', 'go',   'map',  'struct',
    'chan', 'else', 'goto', 'package', 'switch',
    'const', 'fallthrough', 'if',   'range', 'type',
    'continue', 'for',  'import',  'return',  'var',
}

_type_table = {
    UInt64: 'uint64',
    Int64: 'int64',
    UInt32: 'uint32',
    Int32: 'int32',
    Float64: 'float64',
    Float32: 'float32',
    Boolean: 'bool',
    String: 'string',
    # TODO: use GoImportHelper value instead
    Timestamp: 'time.Time',
    Void: 'struct{}',
}


def _rename_if_reserved(s):
    if s in _reserved_keywords:
        return s + '_'
    else:
        return s


def fmt_type(import_helper, data_type, namespace=None, use_interface=False, raw=False):
    data_type, nullable = unwrap_nullable(data_type)
    if is_list_type(data_type):
        if raw and not _needs_base_type(data_type.data_type):
            json_imp = import_helper.id_for_package("encoding/json")
            return "%s.RawMessage" %json_imp
        return '[]%s' % fmt_type(import_helper, data_type.data_type, namespace, use_interface, raw)
    if raw:
        json_imp = import_helper.id_for_package("encoding/json")
        return "%s.RawMessage" % json_imp
    type_name = data_type.name
    if use_interface and _needs_base_type(data_type):
        type_name = 'Is' + type_name
    if is_composite_type(data_type) and namespace is not None and \
            namespace.name != data_type.namespace.name:
        ns_imp = import_helper.id_for_package("github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/%s"%data_type.namespace.name)
        type_name = ns_imp + '.' + type_name
    if use_interface and _needs_base_type(data_type):
        if data_type.__class__ == Timestamp:
            import_helper.id_for_package("time")
        return _type_table.get(data_type.__class__, type_name)
    else:
        if data_type.__class__ == Timestamp:
            import_helper.id_for_package("time")
        return _type_table.get(data_type.__class__, '*' + type_name)


def fmt_var(name, export=True, check_reserved=False):
    s = helpers.fmt_pascal(name) if export else helpers.fmt_camel(name)
    return _rename_if_reserved(s) if check_reserved else s


def _doc_handler(tag, val):
    if tag == 'type':
        return '`{}`'.format(val)
    elif tag == 'route':
        return '`{}`'.format(helpers.fmt_camel(val))
    elif tag == 'link':
        anchor, link = val.rsplit(' ', 1)
        return '`{}` <{}>'.format(anchor, link)
    elif tag == 'val':
        if val == 'null':
            return 'nil'
        else:
            return val
    elif tag == 'field':
        return '`{}`'.format(val)
    else:
        raise RuntimeError('Unknown doc ref tag %r' % tag)


def generate_doc(code_generator, t):
    doc = t.doc
    if doc is None:
        doc = 'has no documentation (yet)'
    doc = code_generator.process_doc(doc, _doc_handler)
    d = '%s : %s' % (fmt_var(t.name), doc)
    if isinstance(t, ApiNamespace):
        d = 'Package %s : %s' % (t.name, doc)
    code_generator.emit_wrapped_text(d, prefix='// ')

    # Generate comment for deprecated routes
    if isinstance(t, ApiRoute):
        if t.deprecated is not None:
            d = 'Deprecated: '
            if t.deprecated.by is not None:
                deprecated_by = t.deprecated.by
                fn = fmt_var(deprecated_by.name)
                if deprecated_by.version != 1:
                    fn += 'V%d' % deprecated_by.version
                d += 'Use `%s` instead' % fn
            code_generator.emit_wrapped_text(d, prefix='// ')


def _needs_base_type(data_type):
    data_type, _ = unwrap_nullable(data_type)
    if is_struct_type(data_type) and data_type.has_enumerated_subtypes():
        return True
    if is_list_type(data_type):
        return _needs_base_type(data_type.data_type)
    return False


def needs_base_type(struct):
    for field in struct.fields:
        if _needs_base_type(field.data_type):
            return True
    return False

GoPackagePath = Text

class GoImportHelper(object):
    """
    A class to keep track of imports to Go files. Every time a package such as
    "dropbox/proto/api_proxy" is needed, the identifier from
    id_for_package("dropbox/proto/api_proxy") should be used instead. This will make sure that no
    two imported packages have conflicting names. Then, the import declarations, including renaming
    imports, can be gotten from emit_import_statements.
    """

    def __init__(self):
        # type: () -> None
        self._package_to_id = dict()  # type: Dict[GoPackagePath, Text]

    def id_for_package(self, package):
        # type: (Text) -> Text
        """
        Gets the id assigned to the package represented by package. If there is no such id, then it
        creates a new id that is guaranteed not to conflict with any other id; if it exists, it will
        deliver the same id.
        """
        if package in self._package_to_id:
            return self._package_to_id[package]
        # The standard identifier for a package in Go is the name of the last directory in the
        # import path. If two different packages would share the same name, at least one of them
        # needs to be renamed. This is done by prepending the directories before the final one, in
        # order, until some free identifier is found. If something has gone wrong and no identifier
        # is assigned this way, we fail an assertion.
        package_list = package.split("/")
        go_package_identifier = u""
        registered = False
        # We iterate through the list in reverse, from deepest to shallowest.
        for package_directory in package_list[-1::-1]:
            if go_package_identifier:
                go_package_identifier = package_directory + u"_" + go_package_identifier
            else:
                # don't have the identifier end in an underscore
                go_package_identifier = package_directory
            if go_package_identifier not in self._package_to_id.values():
                self._package_to_id[package] = go_package_identifier
                registered = True
                break
        if not registered:
            raise NameError(
                "could not find id for package %s go_package_identifier:%s "
                "self._package_to_id:%s"
                % (package, go_package_identifier, self._package_to_id)
            )
        return go_package_identifier

    def emit_import_statements(self, writer):
        # type: (CodeBackend) -> None
        """
        Writes the import statement containing all imports that might be necessary onto the given
        StoneWriter.
        """
        writer.emit("import (")
        with writer.indent():
            for path, ident in sorted(self._package_to_id.items()):
                writer.emit('%s "%s"' % (ident, path))
        writer.emit(")")

    def reset(self):
        """
        Reset the package dict
        """
        self._package_to_id = dict()