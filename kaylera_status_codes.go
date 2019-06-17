package smsproviders

var kaleyraStatusCodeDescription map[string]string = map[string]string{

    "ABSENT-SUB":           "Telecom services not providing service for the particular number. Mobile Subscriber not reachable.",
    "AWAITED-DLR":          "Awaiting delivery report from the operator",
    "BARRED":               "End user has enabled message barring system. Subscriber only accepts messages from Closed User Group [CUG].",
    "BLACKLIST":            "Blacklisted number",
    "DELIVRD":              "SMS successfully delivered.",
    "DNDNUMB":              "DND registered number.",
    "EXPIRED":              "SMS expired after multiple re-try.",
    "FAILED":               "SMS expired due to roaming limitation. / Failed to process the message at operator level.",
    "HANDSET-BUSY":         "Subscriber is in busy condition.",
    "HANDSET-ERR":          "Problem with Handset or handset failed to get the complete message.",
    "INV-NUMBER":           "Invalid number.",
    "INV-TEMPLATE-MATCH":   "Template not matching approved text",
    "INVALID-NUM":          "In case any invalid number present along with the valid numbers.",
    "INVALID-SUB":          "Number does not exist. / Failed to locate the number in HLR database.",
    "MAX-LENGTH":           "In case the message given in any one of the nodes exceeds the maximum of 1000 characters.",
    "MEMEXEC":              "Handset memory full.",
    "MOB-OFF":              "Mobile handset in switched off mode.",
    "NET-ERR":              "Subscriber's operator not supported. / Gateway mobile switching error.",
    "NO-CREDITS":           "Insufficient credits",
    "NO-DLR-OPTR":          "Operator have not acknowledge on status report of the SMS",
    "NOT-OPTIN":            "Not subscribed for opt-in group.",
    "OPTOUT-REJ":           "Optout from subscription.",
    "OUTPUT-REJ":           "Unsubscribed from the group.",
    "REJECTED":             "SMS Rejected as the number is blacklisted by operator.",
    "REJECTED-MULTIPART":   "Validation fail [SMS over 160 characters]",
    "SENDER-ID-NOT-FOUND":  "Sender ID not found",
    "SERIES-BLK":           "Series blocked by the operator.",
    "SERIES-BLOCK":         "Mobile number series blocked.",
    "SERVER-ERR":           "Server error",
    "SNDRID-NOT-ALLOTED":   "Sender ID not allocated",
    "SPAM":                 "Spam SMS",
    "TEMPLATE-NOT-FOUND":   "Template not mapped",
    "TIME-OUT-PROM":        "Time out for promotional SMS.",
    "UNDELIV":              "Failed due to network errors.",
}


func resolveKaleyraStatus(kaleyraStatus string) DeliveryStatus {

    if kaleyraStatus == "AWAITED-DLR" || kaleyraStatus == "NO-DLR-OPTR" {
        return DeliveryStatusAccepted
    } else if kaleyraStatus == "DELIVRD" {
        return DeliveryStatusDelivered
    } else {
        return DeliveryStatusFailed
    }
}


func getKayleraStatusDescription(code string) string {
    return kaleyraStatusCodeDescription[code]
}
