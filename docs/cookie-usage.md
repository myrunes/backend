# Cookie Usage

## What is a cookie?

Here, a short definition from [MDN web docs](https://developer.mozilla.org/en-US/docs/Web/HTTP/Cookies):
> An HTTP cookie (web cookie, browser cookie) is a small piece of data that a server sends to the user's web browser. The browser may store it and send it back with the next request to the same server. Typically, it's used to tell if two requests came from the same browser — keeping a user logged-in, for example. [...]

## How does MYRUNES use cookies?

We are using cookies to store a session key in your browser which authorizes your requests and authenticates your account.  
*In more detail, the session key is a 128 bit base64 encoded string which will be matched against the database to check if there is any user linked to this session.*

## How are local settings saved?

For saving local settings like the acception of the cookie information, we are using [local storage](https://developer.mozilla.org/en-US/docs/Web/API/Window/localStorage). This is not a very secure way of saving private information, but it is verry efficient for local user settings which are not nessecary to be protected like a session key or user ID's.

## Deleting cookies or local storage preferences?

[This](https://support.grammarly.com/hc/en-us/articles/115000090272-How-do-I-clear-cache-and-cookies-) is a very useful guide from gramarly to delete cookies and local data. Just adapt these steps for `myrunes.com`. ;)