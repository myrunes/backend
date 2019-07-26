# REST API

The MYRUNES backend provides a RESTful HTTP JSON API providing 100% of the functionalities of MYRUNES.


## Authenticate

There are two ways to authenticate against the API:

- **API Tokens**  
  In the `MY SETTINGS` page of MYRUNES, you can generate an access token which is a 64 character base64 string used to authenticate against the API.  
  You must pas this token on **each request** as **`Basic`** type token in the **`Authorization`** header. Example:  
  ```
  Authorization: Basic 5lTGAsTFwCKG...
  ```

- **Session Cookies**  
  This method uses session cookies provided by the API which must be saved and delivered on **each request** inside a **`Cookie`** header:  
  ```
  Cookie: __session=tmdV5JEYIB4SzbRmY...
  ```
  To get a session key, request the **[login](#login)** endpoint passing username and password in the JSON body of the request and the server will respond with a **`Set-Cookie`** header containing the session identification after the `__session` key. Keep in mind that you must maintain the expiration of the Cookie because the session will eventually become invalid after this time and must be refreshed. Session are defaultly valid for 2 hours. This value can be extended to a maximum expire duration of 30 days.


## Body Content Type

All request bodies which are sent to the API or returned by the API are using content type **`application/json`**.


## Request Parameters

Optional request parameters are styled *`like this`*  and non-optionals are styled `like this` in the documentation. There are different types how parameters must be passed. 

- Either as `Path` variable as part of the request resource like `/api/pages/`**`87128927213891273`** or
- as `URL Query` like `/api/pages`**`?sort_by=date`** or
- as `Body` parameter like  
  ```json
  {
    "username": "zekro"
  }
  ```

## Error Responses

The API uses the standart HTTP/1.1 status codes like defined in [RFC 2616 - section 10](https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html).

Also, every error response contains a body with the status `code` and an error `message` as description of the error.

```json
{
    "code": 429,
    "message": "You are being rate limited"
}
```

## Rate Limiting

The API is rate limited by a per-connection and per-endpoint [token bucket](https://en.wikipedia.org/wiki/Token_bucket) limiter system. Also, there is a global limiter across all endpoints per-connection.

For each endpoint, you will have a maximum ammount of tokens you can use for requests. Each request, one token will be consumed. Each time a specified ammount of time elapses, a new token will be added to your bucket.

You can check your current rate limit status by examining the passed headers
- `X-Ratelimit-Limit`  
   which displays the total ammount of maximum token you can have,
- `X-Ratelimit-Remaining`  
   which presents the ammount of tokens you can stil use,
- `X-Ratelimit-Reset`  
   gives the UNIX time stamp (seconds) when you are able to request again after consumption of all tokens

---

## API Objects

There are different types of API objects which will be returned by the API.

### User Object

> A MYRUNES registered user account.

| Key | Type |  Description |
|-----|------|--------------|
| `uid` | string | Unique user ID in form of a [snowflake](https://developer.twitter.com/en/docs/basics/twitter-ids.html) like object |
| `username` | string | The unique user name of the user |
| `displayname` | string | The display name of the user (may not be unique) |
| `lastlogin` | string | Time of last successful login |
| `created` | string | Time of account creation |
| `favorites` | List\<string\> | List of favorited champion IDs |

```json
{
  "uid": "1154685560976457728",
  "username": "zekro",
  "displayname": "zekro der Echte",
  "lastlogin": "2019-07-26T09:31:05.62Z",
  "created": "2019-07-26T09:31:04.993Z",
  "favorites": [
    "kindred",
    "pyke",
    "lux"
  ]
}
```

### Page Object

> A rune page.

Rune pages consists of following sub-objects.  
A full list featuring all available trees, runes and perks you can get py requesting the [`resources`](#resources) paths.

**Primary Tree Object**

> The primary tree of a rune page.

```json
{
  "tree": "domination",
  "rows": [
    "electrocute",
    "cheap-shot",
    "zombie-ward",
    "ravenous-hunter"
  ]
}
```

**Secondary Tree Object**

>  The secondary tree of a rune page.

```json
{
  "tree": "precission",
  "rows": [
    "legend-bloodline",
    "cut-down"
  ]
}
```

**Perks Object**

> The perks collection of a rune page.

```json
{
  "rows": [
    "diamond",
    "diamond",
    "heart"
  ]
}
```

The actual page object is built like follwing:

| Key | Type |  Description |
|-----|------|--------------|
| `uid` | string | Unique page ID in form of a [snowflake](https://developer.twitter.com/en/docs/basics/twitter-ids.html) like object |
| `owner` | string | The user UID of the owner/creator of the page |
| `title` | string | The title of the page |
| `created` | string | The date of creation of the page |
| `edited` | string | The date of the last modification of the page |
| `champions` | List\<string\> | List of champion IDs the page is linked to |
| `primary` | Primary Tree Object | |
| `secondary` | Secondary Tree Object | |
| `perks` | Perks Object | |

```json
{
  "uid": "1136539895017013248",
  "owner": "1136250237250584576",
  "title": "asdasd",
  "created": "2019-06-06T07:46:41.517Z",
  "edited": "2019-06-11T13:02:53.124Z",
  "champions": [
    "lux"
  ],
  "primary": { Primary Page Object },
  "secondary": { Secondary Page Object },
  "perks": { Perks Object }
}
```

### Share Object

> A representation of data of a shared rune page.

| Key | Type |  Description |
|-----|------|--------------|
| `uid` | string | Unique share ID in form of a [snowflake](https://developer.twitter.com/en/docs/basics/twitter-ids.html) like object |
| `ident` | string | The unique random identifier in format of a 8 character long base64 string used to request the shared page and represent in the share link |
| `owner` | string | The unique ID of the owner of the shared page |
| `page` | string | The unique ID of the shared page |
| `created` | string | Date of the creation of the share |
| `maxaccesses` | number | Maximum ammount of accesses available |
| `accesses` | number | Ammount of accesses until now |
| `expires` | string | The date of expiration. This will alway be a valid parsable value even though expiration is not set, this will be a time very far in the future |
| `lastaccess` | string | Date of the last access |
| `page` | Page Object | The object of the shared page |
| `user` | User Object | The object of the user opened the share |

```json
{
  "share": {
    "uid": "1145956056344100864",
    "ident": "z9qrkRQ=",
    "owner": "1136250237250584576",
    "page": "1136961585131847680",
    "created": "2019-07-02T07:23:09.323Z",
    "maxaccesses": 4,
    "expires": "2119-06-08T07:23:09.323Z",
    "accesses": 1,
    "lastaccess": "2019-07-02T07:26:52.875Z"
  },
  "page": { Page Object },
  "user": { User Object }
}
```

### Session Object

> Information around a session key bound to a user to authenticate their requests.

| Key | Type |  Description |
|-----|------|--------------|
| `sessionid` | string | Unique session ID in form of a [snowflake](https://developer.twitter.com/en/docs/basics/twitter-ids.html) like object |
| `key` | string | A pseudo-representation of the session key showing the first and last 3 characters of the key |
| `uid` | string | Unique ID of the user bound to this session | 
| `expires` | string | Date when the session key turns invalid |
| `lastaccess` | string | Date of the last authentication using this session key |
| `lastaccessip` | string | The remote address of the last authenticated request using this session kes |

```json
{
  "sessionid": "1154652075600297984",
  "key": "f07...l8=",
  "uid": "1136250237250584576",
  "expires": "2019-08-25T07:18:01.879Z",
  "lastaccess": "2019-07-26T09:30:52.085Z",
  "lastaccessip": "127.0.0.1"
}
```

### API Token Object

> Information around an API token bound to a user to authenticate their requests.

| Key | Type |  Description |
|-----|------|--------------|
| `userid` | string | The unique ID of the user bound to this token |
| `token` | string | The API token secret |
| `created` | string | Date the API token was generated |

---

## Endpoints

### Authentication

#### Login

> `POST /api/login`

**Parameters**

| Name | Type | Via | Default | Description |
|------|------|-----|---------|-------------|
| `username` | string | Body | | The username of the account |
| `password` | string | Body | | The password of the given user |
| *`remember`* | boolean | Body | `false` | Sessions defaultly expire after 2 hours. Setting this to true, this duration will be expanded to 30 days. |

**Response**

```
HTTP/1.1 200 OK
Content-Length: 36
Content-Type: application/json
Date: Fri, 26 Jul 2019 11:04:02 GMT
Server: MYRUNES v.DEBUG_BUILD
X-Ratelimit-Limit: 50
X-Ratelimit-Remaining: 48
X-Ratelimit-Reset: 0
```
```json
{
  "code": 200,
  "message": "ok"
}
```

#### Logout

> `POST /api/logout`

**Parameters**

*No parameters necessary.*

**Response**

```
HTTP/1.1 200 OK
Content-Length: 36
Content-Type: application/json
Date: Fri, 26 Jul 2019 11:06:43 GMT
Server: MYRUNES v.DEBUG_BUILD
X-Ratelimit-Limit: 50
X-Ratelimit-Remaining: 49
X-Ratelimit-Reset: 0
```
```json
{
  "code": 200,
  "message": "ok"
}
```

### Users

#### Get Self User

> `GET /api/users/me`

**Parameters**

*No parameters necessary.*

**Response**

```
HTTP/1.1 200 OK
Content-Length: 199
Content-Type: application/json
Date: Fri, 26 Jul 2019 11:21:07 GMT
Server: MYRUNES v.DEBUG_BUILD
X-Ratelimit-Limit: 50
X-Ratelimit-Remaining: 49
X-Ratelimit-Reset: 0
```
```json
{ User Object }
```

#### Check User Name

> `GET /api/users/:USERNAME`

*This endpoint is concipated for checking the availability of a username on registration, not to gather user information from another account which is not possible yet over the API.*  
*If the given username is unused, a 404 Not Found response will be returned which then should be interpreted as success or available.*

**Parameters**

| Name | Type | Via | Default | Description |
|------|------|-----|---------|-------------|
| `USERNAME` | string | Path | | The username to be checked |

**Response**

```
HTTP/1.1 200 OK
Content-Length: 36
Content-Type: application/json
Date: Fri, 26 Jul 2019 11:08:07 GMT
Server: MYRUNES v.DEBUG_BUILD
X-Ratelimit-Limit: 50
X-Ratelimit-Remaining: 49
X-Ratelimit-Reset: 0
```
```json
{
  "code": 200,
  "message": "ok"
}
```

#### Create User

> `POST /api/users`

**Parameters**

| Name | Type | Via | Default | Description |
|------|------|-----|---------|-------------|
| `username` | string | Body | | The username of the account |
| `password` | string | Body | | The password of the given user |
| *`remember`* | boolean | Body | `false` | Sessions defaultly expire after 2 hours. Setting this to true, this duration will be expanded to 30 days. |

**Response**

```
HTTP/1.1 201 Created
Content-Length: 210
Content-Type: application/json
Date: Fri, 26 Jul 2019 11:28:40 GMT
Server: MYRUNES v.DEBUG_BUILD
X-Ratelimit-Limit: 1
X-Ratelimit-Remaining: 0
X-Ratelimit-Reset: 1564140534
```
```json
{ User Object }
```

#### Update Self User

> `POST /api/users/me`

**Parameters**

| Name | Type | Via | Default | Description |
|------|------|-----|---------|-------------|
| `currpassword` | string | Body | | The current password of the users account |
| *`newpassword`* | string | Body | | A new password which will replace the current one |
| *`displayname`* | string | Body | | A new display name |
| *`username`* | string | Body | | A new user name |

**Response**

```
HTTP/1.1 200 OK
Content-Length: 36
Content-Type: application/json
Date: Fri, 26 Jul 2019 11:34:22 GMT
Server: MYRUNES v.DEBUG_BUILD
X-Ratelimit-Limit: 50
X-Ratelimit-Remaining: 48
X-Ratelimit-Reset: 0
```
```json
{
  "code": 200,
  "message": "ok"
}
```

#### Delete Self User

> `DELETE /api/users/me`

**Parameters**

| Name | Type | Via | Default | Description |
|------|------|-----|---------|-------------|
| `currpassword` | string | Body | | The current password of the users account |

**Response**

```
HTTP/1.1 200 OK
Content-Length: 36
Content-Type: application/json
Date: Fri, 26 Jul 2019 11:39:15 GMT
Server: MYRUNES v.DEBUG_BUILD
X-Ratelimit-Limit: 50
X-Ratelimit-Remaining: 48
X-Ratelimit-Reset: 0
```
```json
{
  "code": 200,
  "message": "ok"
}
```

