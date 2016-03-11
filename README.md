Package ctxmiddleware is an example of using custom context handlers and
middleware. It is largely based on the final option in this article:
https://joeshaw.org/net-context-and-http-handler/.
I have added a middleware chaining method and an associated ContextMW type.

The idea is that in the case of this proposal being accepted,
https://github.com/golang/go/issues/14660,
where the context package would be in the standard library and implemented
on http.Request, that these same handlers and middleware could be used after
a brief refactoring to remove the custom types.
