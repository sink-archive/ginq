# GINQ

**<u>G</u>o <u>IN</u>tegrated <u>Q</u>uery**



[![`GPL-3.0-or-later`](https://img.shields.io/badge/license-GPL--3.0--or--later-blue)](https://github.com/yellowsink/ginq/blob/master/LICENSE.md)

GINQ ports the amazing LINQ features of .NET to Golang. Note that this library only deals with the `.Where()`, `.Select()`, and similar extension methods. The awesome `(from x in list select x.y where x.z` query syntax would require special language support that a library cannot provide.

Read about LINQ [here!](https://docs.microsoft.com/en-us/dotnet/standard/linq)