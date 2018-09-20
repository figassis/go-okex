package okex

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/shopspring/decimal"
)

const (
	URL = "https://www.okex.com/api/v1/"
)

func New(key, secret string) Client {
	return Client{APIKey: key, SecretKey: secret}
}

func (c Client) GetUserInfo() (UserInfo, error) {
	var params = url.Values{}

	params.Set("api_key", c.APIKey)

	sendAuthenticatedReq(&params, c.SecretKey)

	return okexUserInfo(params)

}

func (c Client) Trade(symbol string, orderType string, price decimal.Decimal, amount decimal.Decimal) (Receipt, error) {
	var params = url.Values{}

	params.Set("api_key", c.APIKey)
	params.Set("symbol", symbol)
	params.Set("type", orderType)
	params.Set("price", price.String())
	params.Set("amount", amount.String())

	sendAuthenticatedReq(&params, c.SecretKey)

	return okexTrade(params)

}

func (c Client) cancelOrder(symbol string, orderId string) (Receipt, error) {
	var params = url.Values{}

	params.Set("api_key", c.APIKey)
	params.Set("symbol", symbol)
	params.Set("order_id", orderId)

	sendAuthenticatedReq(&params, c.SecretKey)

	return okexCancel(params)

}

func (c Client) getOrderInfo(symbol string, orderId string) (OrderInfo, error) {
	var params = url.Values{}

	params.Set("api_key", c.APIKey)
	params.Set("symbol", symbol)
	params.Set("order_id", orderId)

	sendAuthenticatedReq(&params, c.SecretKey)

	return okexOrderInfo(params)

}

func (c Client) getOrderHistory(symbol string, status string, currentPage string, pageLength string) (OrderHistory, error) {
	var params = url.Values{}

	params.Set("api_key", c.APIKey)
	params.Set("symbol", symbol)
	params.Set("status", status)
	params.Set("current_page", currentPage)
	params.Set("page_length", pageLength)

	sendAuthenticatedReq(&params, c.SecretKey)

	return okexOrderHistory(params)

}

///// helper functions

func reqPost(url string, c string) string {
	request := gorequest.New().Timeout(time.Second * 5)
	resp, body, err := request.Post(url).
		Set("contentType", "application/x-www-form-urlencoded").
		Send(c).
		End()

	if err != nil && resp.StatusCode != 200 {
		panic(err)
	}

	return body

}

func sendAuthenticatedReq(postForm *url.Values, secretKey string) error {
	payload := postForm.Encode()
	payload = payload + "&secret_key=" + secretKey
	payload2, _ := url.QueryUnescape(payload)

	sign, err := GetParamMD5Sign(secretKey, payload2)
	if err != nil {
		return err
	}

	postForm.Set("sign", strings.ToUpper(sign))
	return nil
}

func GetParamMD5Sign(secret, params string) (string, error) {
	hash := md5.New()
	_, err := hash.Write([]byte(params))

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func okexUserInfo(params url.Values) (UserInfo, error) {
	var objmap = UserInfo{}
	var error = ErrorCode{}

	resp := reqPost(URL+"userinfo.do", params.Encode())

	if err := json.Unmarshal([]byte(resp), &error); err != nil {
	}
	if err := json.Unmarshal([]byte(resp), &objmap); err != nil {
	}

	return objmap, getErrorMessage(error.ErrorCode)

}

func okexTrade(params url.Values) (Receipt, error) {
	var obj = Receipt{}
	var error = ErrorCode{}

	resp := reqPost(URL+"trade.do", params.Encode())

	if err := json.Unmarshal([]byte(resp), &error); err != nil {
	}
	if err := json.Unmarshal([]byte(resp), &obj); err != nil {
	}

	return obj, getErrorMessage(error.ErrorCode)

}

func okexCancel(params url.Values) (Receipt, error) {
	var obj = Receipt{}
	var error = ErrorCode{}

	resp := reqPost(URL+"cancel_order.do", params.Encode())

	if err := json.Unmarshal([]byte(resp), &error); err != nil {
	}
	if err := json.Unmarshal([]byte(resp), &obj); err != nil {
	}

	return obj, getErrorMessage(error.ErrorCode)

}

func okexOrderInfo(params url.Values) (OrderInfo, error) {
	var obj = OrderInfo{}
	var error = ErrorCode{}

	resp := reqPost(URL+"order_info.do", params.Encode())

	if err := json.Unmarshal([]byte(resp), &error); err != nil {
	}
	if err := json.Unmarshal([]byte(resp), &obj); err != nil {
	}

	return obj, getErrorMessage(error.ErrorCode)

}

func okexOrderHistory(params url.Values) (OrderHistory, error) {
	var obj = OrderHistory{}
	var error = ErrorCode{}

	resp := reqPost(URL+"order_history.do", params.Encode())

	if err := json.Unmarshal([]byte(resp), &error); err != nil {
	}
	if err := json.Unmarshal([]byte(resp), &obj); err != nil {
	}

	return obj, getErrorMessage(error.ErrorCode)

}

func GetSlash() string {
	if runtime.GOOS == "windows" {
		return "\\"
	}
	return "/"
}

func getErrorMessage(_input int) error {
	var input = strconv.Itoa(_input)

	var ErrorCodes = map[string]error{
		//Spot Errors
		"10000": errors.New("Required field, can not be null"),
		"10001": errors.New("Request frequency too high to exceed the limit allowed"),
		"10002": errors.New("System error"),
		"10004": errors.New("Request failed"),
		"10005": errors.New("'SecretKey' does not exist"),
		"10006": errors.New("'Api_key' does not exist"),
		"10007": errors.New("Signature does not match"),
		"10008": errors.New("Illegal parameter"),
		"10009": errors.New("Order does not exist"),
		"10010": errors.New("Insufficient funds"),
		"10011": errors.New("Amount too low"),
		"10012": errors.New("Only btc_usd ltc_usd supported"),
		"10013": errors.New("Only support https request"),
		"10014": errors.New("Order price must be between 0 and 1,000,000"),
		"10015": errors.New("Order price differs from current market price too much"),
		"10016": errors.New("Insufficient coins balance"),
		"10017": errors.New("API authorization error"),
		"10018": errors.New("borrow amount less than lower limit [usd:100,btc:0.1,ltc:1]"),
		"10019": errors.New("loan agreement not checked"),
		"10020": errors.New("rate cannot exceed 1%"),
		"10021": errors.New("rate cannot less than 0.01%"),
		"10023": errors.New("fail to get latest ticker"),
		"10024": errors.New("balance not sufficient"),
		"10025": errors.New("quota is full, cannot borrow temporarily"),
		"10026": errors.New("Loan (including reserved loan) and margin cannot be withdrawn"),
		"10027": errors.New("Cannot withdraw within 24 hrs of authentication information modification"),
		"10028": errors.New("Withdrawal amount exceeds daily limit"),
		"10029": errors.New("Account has unpaid loan, please cancel/pay off the loan before withdraw"),
		"10031": errors.New("Deposits can only be withdrawn after 6 confirmations"),
		"10032": errors.New("Please enabled phone/google authenticator"),
		"10033": errors.New("Fee higher than maximum network transaction fee"),
		"10034": errors.New("Fee lower than minimum network transaction fee"),
		"10035": errors.New("Insufficient BTC/LTC"),
		"10036": errors.New("Withdrawal amount too low"),
		"10037": errors.New("Trade password not set"),
		"10040": errors.New("Withdrawal cancellation fails"),
		"10041": errors.New("Withdrawal address not exsit or approved"),
		"10042": errors.New("Admin password error"),
		"10043": errors.New("Account equity error, withdrawal failure"),
		"10044": errors.New("fail to cancel borrowing order"),
		"10047": errors.New("this function is disabled for sub-account"),
		"10048": errors.New("withdrawal information does not exist"),
		"10049": errors.New("User can not have more than 50 unfilled small orders (amount<0.15BTC)"),
		"10050": errors.New("can't cancel more than once"),
		"10051": errors.New("order completed transaction"),
		"10052": errors.New("not allowed to withdraw"),
		"10064": errors.New("after a USD deposit, that portion of assets will not be withdrawable for the next 48 hours"),
		"10100": errors.New("User account frozen"),
		"10101": errors.New("order type is wrong"),
		"10102": errors.New("incorrect ID"),
		"10103": errors.New("the private otc order's key incorrect"),
		"10216": errors.New("Non-available API"),
		"1002":  errors.New("The transaction amount exceed the balance"),
		"1003":  errors.New("The transaction amount is less than the minimum requirement"),
		"1004":  errors.New("The transaction amount is less than 0"),
		"1007":  errors.New("No trading market information"),
		"1008":  errors.New("No latest market information"),
		"1009":  errors.New("No order"),
		"1010":  errors.New("Different user of the cancelled order and the original order"),
		"1011":  errors.New("No documented user"),
		"1013":  errors.New("No order type"),
		"1014":  errors.New("No login"),
		"1015":  errors.New("No market depth information"),
		"1017":  errors.New("Date error"),
		"1018":  errors.New("Order failed"),
		"1019":  errors.New("Undo order failed"),
		"1024":  errors.New("Currency does not exist"),
		"1025":  errors.New("No chart type"),
		"1026":  errors.New("No base currency quantity"),
		"1027":  errors.New("Incorrect parameter may exceeded limits"),
		"1028":  errors.New("Reserved decimal failed"),
		"1029":  errors.New("Preparing"),
		"1030":  errors.New("Account has margin and futures, transactions can not be processed"),
		"1031":  errors.New("Insufficient Transferring Balance"),
		"1032":  errors.New("Transferring Not Allowed"),
		"1035":  errors.New("Password incorrect"),
		"1036":  errors.New("Google Verification code Invalid"),
		"1037":  errors.New("Google Verification code incorrect"),
		"1038":  errors.New("Google Verification replicated"),
		"1039":  errors.New("Message Verification Input exceed the limit"),
		"1040":  errors.New("Message Verification invalid"),
		"1041":  errors.New("Message Verification incorrect"),
		"1042":  errors.New("Wrong Google Verification Input exceed the limit"),
		"1043":  errors.New("Login password cannot be same as the trading password"),
		"1044":  errors.New("Old password incorrect"),
		"1045":  errors.New("2nd Verification Needed"),
		"1046":  errors.New("Please input old password"),
		"1048":  errors.New("Account Blocked"),
		"1201":  errors.New("Account Deleted at 00: 00"),
		"1202":  errors.New("Account Not Exist"),
		"1203":  errors.New("Insufficient Balance"),
		"1204":  errors.New("Invalid currency"),
		"1205":  errors.New("Invalid Account"),
		"1206":  errors.New("Cash Withdrawal Blocked"),
		"1207":  errors.New("Transfer Not Support"),
		"1208":  errors.New("No designated account"),
		"1209":  errors.New("Invalid api"),
		"1216":  errors.New("Market order temporarily suspended. Please send limit order"),
		"1217":  errors.New("Order was sent at ±5% of the current market price. Please resend"),
		"1218":  errors.New("Place order failed. Please try again later"),
		// Errors for both
		"HTTP ERROR CODE 403": errors.New("Too many requests, IP is shielded"),
		"Request Timed Out":   errors.New("Too many requests, IP is shielded"),
		// contract errors
		"405":   errors.New("method not allowed"),
		"20001": errors.New("User does not exist"),
		"20002": errors.New("Account frozen"),
		"20003": errors.New("Account frozen due to liquidation"),
		"20004": errors.New("Contract account frozen"),
		"20005": errors.New("User contract account does not exist"),
		"20006": errors.New("Required field missing"),
		"20007": errors.New("Illegal parameter"),
		"20008": errors.New("Contract account balance is too low"),
		"20009": errors.New("Contract status error"),
		"20010": errors.New("Risk rate ratio does not exist"),
		"20011": errors.New("Risk rate lower than 90%/80% before opening BTC position with 10x/20x leverage. or risk rate lower than 80%/60% before opening LTC position with 10x/20x leverage"),
		"20012": errors.New("Risk rate lower than 90%/80% after opening BTC position with 10x/20x leverage. or risk rate lower than 80%/60% after opening LTC position with 10x/20x leverage"),
		"20013": errors.New("Temporally no counter party price"),
		"20014": errors.New("System error"),
		"20015": errors.New("Order does not exist"),
		"20016": errors.New("Close amount bigger than your open positions"),
		"20017": errors.New("Not authorized/illegal operation"),
		"20018": errors.New("Order price cannot be more than 103% or less than 97% of the previous minute price"),
		"20019": errors.New("IP restricted from accessing the resource"),
		"20020": errors.New("c.SecretKey does not exist"),
		"20021": errors.New("Index information does not exist"),
		"20022": errors.New("Wrong API interface (Cross margin mode shall call cross margin API, fixed margin mode shall call fixed margin API)"),
		"20023": errors.New("Account in fixed-margin mode"),
		"20024": errors.New("Signature does not match"),
		"20025": errors.New("Leverage rate error"),
		"20026": errors.New("API Permission Error"),
		"20027": errors.New("no transaction record"),
		"20028": errors.New("no such contract"),
		"20029": errors.New("Amount is large than available funds"),
		"20030": errors.New("Account still has debts"),
		"20038": errors.New("Due to regulation, this function is not available in the country/region your currently reside in"),
		"20049": errors.New("Request frequency too high"),
	}

	return ErrorCodes[input]

}
