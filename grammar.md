[Stmt]        | ([Assignment] | [Expr] | [Exit] | [IfEq] | [Switch])
[IfEq]        | 'ifeq' [name] [Int] ([Stmt]?)
[Switch]      | 'switch' [name] ([case]?)
[case]        | 'case' [name] [Int] ([Stmt]?)
[Assignment]  | ([name] [Int])
[Expr]        | ([Push] | [Pop] | [ChOut])
[Push]        | 'push' [Int]
[Pop]         | 'pop'
[ChOut]       | 'chout'
[Name]        | characters
[Int]         | digits
[Exit]        |'Exit'

===== Example code



