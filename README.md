


# Iracema
### Syntax

```iracema
object Lang {
  fun init(name) {
    @name = name
  }

  fun name {
    return @name
  }
}

l = Lang.new("Iracema")
puts(l.name)

```
### Arithmetic Operators
Following table shows all the arithmetic operators supported by Iracema. Assume variable **A** holds 10 and variable **B** holds 20 then
| Operator | Description | Example |
| --- | --- | --- |
| +   | Adds two operands | A + B will give 30 |
| -   | Subtracts second operand from the first | A - B will give -10 |
| *   | Multiply both operands | A * B will give 200 |
| /   | Divide numerator by de-numerator | B / A will give 2 |
| -   | Unary - operator acts as negation | -A will give -10 |

### Keywords
The following list shows a few of the reserved words in Iracema. These reserved words may not be used as constants or variables or any other identifier names.

<table>
<body>
 <tr>
    <td>object</td>
    <td>fun</td>
    <td>catch</td>
    <td>return</td>
    <td>stop</td>
     <td>next</td>
  </tr>
  <tr>
    <td>if</td>
    <td>else</td>
    <td>while</td>
    <td>true</td>
    <td>false</td>
     <td>nil</td>
  </tr>
</tbody>
</table>



### Installation from source

1. Verify that you have Go 1.17+ installed

   ```sh
   $ go version
   ```

   If `go` is not installed, follow instructions on [the Go website](https://golang.org/doc/install).

2. Clone this repository

   ```sh
   $ https://github.com/alissonbrunosa/iracema
   $ cd iracema
   ```

3. Build

   ```sh
   $  go build cmd/iracema/main.go -o iracema
   ```
