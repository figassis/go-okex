**Simple REST api for Okex exchange**

Install by typing: 

```sh
go get -u github.com/figassis/go-okex
```

Currently supported api calls:

+ Get User Info : returns user info (funds)
+ Trade: place a trade
+ Cancel order: cancel a standing order
+ Get Order Info: returns information on an order based on the order Id
+ Get Order History: returns information on orders from the last two days for a given coin

Implementation based on the official Okex api [documentation](https://github.com/okcoin-okex/API-docs-OKEx.com/blob/master/API-For-Spot-EN/REST%20API%20for%20SPOT.md).

Using the [GoRequest](https://github.com/parnurzeal/gorequest) http client to send calls and some functions have been copied from [GoCryptoTrader](https://github.com/thrasher-/gocryptotrader). Thank you!

Have fun :)

