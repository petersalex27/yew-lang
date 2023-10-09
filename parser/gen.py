import sys

name = sys.argv[1]

out = ( 
  f'package parser\n' +
  f'\n' +
  f'import (\n' +
  f'\t"github.com/petersalex27/yew-packages/parser"\n' +
  f'\t"github.com/petersalex27/yew-packages/parser/ast"\n' +
  f'\titoken "github.com/petersalex27/yew-packages/token"\n' +
  f')\n' +
  f'\n'
  f'type {name}Node struct' + '{\n' +
  f'\n' +
  '}\n' +
  f'\n' +
  f'func (n {name}Node) Equals(a ast.Ast) bool ' + '{\n' +
  f'\tn2, ok := a.({name}Node)\n' +
  f'\tif !ok ' + '{\n' +
  f'\t\treturn false\n' +
  '\t}\n' +
  f'\t//TODO: finish\n' +
  '}\n' +
  f'\n' +
  f'func (n {name}Node) NodeType() ast.Type ' + '{ return ' + f'{name}' + ' }\n' +
  f'\n' +
  f'func (n {name}Node) InOrderTraversal(f func(itoken.Token)) ' + '{\n' +
  f'\tpanic("implement: ({name}Node) InOrderTraversal(func(itoken.Token))")\n' +
  '}'
)

print(out)