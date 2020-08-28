module github.com/KlyuchnikovV/cui

go 1.14

require (
	github.com/KlyuchnikovV/lines_buffer v0.0.0-20200805140229-a63fde50d25b
	github.com/KlyuchnikovV/termin v0.0.0-20200805143708-7655119a7710
)

replace (
	github.com/KlyuchnikovV/lines_buffer => ../lines_buffer
	github.com/KlyuchnikovV/termin => ../termin
)
