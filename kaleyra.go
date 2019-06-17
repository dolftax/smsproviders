package smsproviders

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "time"
    "io/ioutil"
    "net/http"
)

/////////////////////////////////////
// Kaylera specific configurations //
/////////////////////////////////////

const (
    kaleyraSlug string = "KALEYRA"
    kaleyraAPIUrl string = "https://api-alerts.kaleyra.com/v4"
    kaleyraTimeFormat string = "2006-01-02 15:04:05-0700"
)

type KayleraSMSData struct {
    Id              string `json:"id"`
    CustomId        string `json:"customid"`
    Custom          string `json:"custom"`
    Mobile          string `json:"mobile"`
    Status          string `json:"status"`
    DlrTimeString   string `json:"dlrtime"`
}

type KaleyraAPIResponse struct {
    Status          string `json:"status"`
    Data            []*KayleraSMSData `json:"data"`
    Code            string `json:"code"`
    Message         string `json:"message"`
}


func parseKayleraTime(timeString string) (time.Time, error) {

    var kaleyraTime time.Time
    // Kaleyra sends time as string in the format
    // `2006-01-02 15:04:05`. This time is in IST but does not
    // contain any time zone related data. Hence we need to add
    // "+0530" to the time string
    timeString += "+0530"
    kaleyraTime, err := time.Parse(kaleyraTimeFormat, timeString)
    if err != nil {
        kaleyraTime = kaleyraTime.UTC()
    }

    return kaleyraTime, err
}

func getReportFromApiResponse(response *http.Response) (*SMSStatusReport, error) {

    // If the http response does not have a 2XX response, the 
    // AttemptResponse will have errorCode as "HTTP:<StatusCode>"
    // and errorMessage as the response body.
    if !(response.StatusCode >= 200 && response.StatusCode < 300) {
        body, _ := ioutil.ReadAll(response.Body)
        log.Printf("[KALEYRA] SMS Atempt HTTP Response [%v]:\nBody: %v",
                   response.StatusCode, string(body))

        return nil, fmt.Errorf("HTTP Respone Status Code: %vStatus")
    }

    // Parse http response to KaleyraAPIResponse object
    // using JSON decoder
    var apiResponse KaleyraAPIResponse
    decoder := json.NewDecoder(response.Body)
    err := decoder.Decode(&apiResponse)
    if err != nil {
        return nil, err
    }

    report := &SMSStatusReport{
        Provider: kaleyraSlug,
    }

    if apiResponse.Status == "OK" {

        data := apiResponse.Data[0]
        report.MessageId = data.Id
        report.CustomId = data.CustomId
        report.ProviderStatus = data.Status
        report.DeliveryStatus = resolveKaleyraStatus(data.Status)
        report.Description = getKayleraStatusDescription(data.Status)

        dlrTime, _ := parseKayleraTime(data.DlrTimeString)
        report.DeliveryTime = &dlrTime

    } else if apiResponse.Status == "ERROR" {

        report.ProviderStatus = apiResponse.Code
        report.DeliveryStatus = resolveKaleyraStatus(apiResponse.Code)
        report.Description = getKayleraStatusDescription(apiResponse.Code)
    
    } else {
        return nil, fmt.Errorf(
            "[KALEYRA] Unknown STATUS'%v'. Expected OK/ERROR.", apiResponse.Status)
    }

    return report, nil
}


type Kaleyra struct {
    APIKey        string
    SenderId      string
    CallbackUrl   string
    Client        *http.Client
}

func NewKaleyraClient() Kaleyra {

    client := Kaleyra{
        APIKey: os.Getenv("KALEYRA_API_KEY"),
        SenderId: os.Getenv("KALEYRA_SENDER_ID"),
        CallbackUrl: os.Getenv("KALEYRA_CALLBACK_URL"),
        Client: &http.Client{},
    }

    return client
}

func (provider Kaleyra) SendSMS(recipient, message, customId string) (*SMSStatusReport, error) {

    // Create Http Request

    request, _ := http.NewRequest("GET", kaleyraAPIUrl, nil)

    if customId == "" {
        customId = "1"
    }

    query := request.URL.Query()
    query.Add("method", "sms")
    query.Add("api_key", provider.APIKey)
    query.Add("sender", provider.SenderId)
    query.Add("to", recipient)
    query.Add("message", message)
    query.Add("custom", customId)

    // Kaleyra's callback request sends `status` (ProviderStatus) and
    // `msgid` (CustomId) in the request by default. In order to get
    // provider's message_id and delivery time, we need to add `dlrtime`
    // and `sid` to get them in the callback params:
    // Ref: https://promo.solutionsinfini.com/readme/4.0/send-sms-xml 
    callbackUrl := provider.CallbackUrl + "?dlrtime={delivered}&id={sid}"
    query.Add("dlrurl", callbackUrl)

    request.URL.RawQuery = query.Encode()

    // Send the request
    response, _ := provider.Client.Do(request)
    return getReportFromApiResponse(response)    
}

func (provider Kaleyra) GetStatusReport(messageId string) (*SMSStatusReport, error) {

    // Create Http Request
    request, _ := http.NewRequest("GET", kaleyraAPIUrl, nil)
    
    query := request.URL.Query()
    query.Add("method", "sms.status")
    query.Add("api_key", provider.APIKey)
    query.Add("id", messageId)

    request.URL.RawQuery = query.Encode()

    // Send Http Request
    response, _ := provider.Client.Do(request)
    return getReportFromApiResponse(response)
}

func (provider Kaleyra) CallbackToReport(request *http.Request) (*SMSStatusReport, error) {

    data := KayleraSMSData{}
    params := request.URL.Query()

    mIdString, ok := params["sid"]
    if ok && len(mIdString) >= 1 {
        data.Id = mIdString[0]
    } else {
        return nil, fmt.Errorf("[KALEYRA] Unable to get 'id' in callback request")
    }

    providerStatus, ok := params["status"]
    if ok && len(providerStatus) >= 1 {
        data.Status = providerStatus[0]
    } else {
        return nil, fmt.Errorf("[KALEYRA] Unable to get 'status' in callback request")
    }

    customId, ok := params["msgid"]
    if ok && len(customId) >= 1 {
        data.CustomId = customId[0]
    }

    report := &SMSStatusReport{
        Provider: kaleyraSlug,
        MessageId: data.Id,
        CustomId: data.CustomId,
        ProviderStatus: data.Status,
        DeliveryStatus: resolveKaleyraStatus(data.Status),
        Description: getKayleraStatusDescription(data.Status),
    }

    dlrtimeString, ok := params["dlrtime"]
    if ok && len(dlrtimeString) >= 1{
        data.DlrTimeString = dlrtimeString[0]
        deliveryTime, _ := parseKayleraTime(data.DlrTimeString)
        if !deliveryTime.IsZero() {
            report.DeliveryTime = &deliveryTime
        }
    }

    return report, nil
}
