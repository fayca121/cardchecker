package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type cardNumberRequest struct {
	CardNumber string `json:"card_number"`
}

type cardNumberResponse struct {
	cardNumberRequest
	Valid   bool   `json:"valid"`
	Network string `json:"network,omitempty"`
}

func CheckCardNumber(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	var req cardNumberRequest

	w.Header().Set("Content-Type", "application/json")

	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			panic(err)
		}
	}(r.Body)

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rp := strings.NewReplacer(" ", "", "-", "")

	cardNumber := rp.Replace(req.CardNumber)

	isNumeric := regexp.MustCompile(`^[0-9]+$`).MatchString(cardNumber)

	if len(cardNumber) < 8 || len(cardNumber) > 19 || !isNumeric {
		http.Error(w, "Invalid card number", http.StatusBadRequest)
		return
	}

	resp := cardNumberResponse{
		cardNumberRequest: req,
		Valid:             isValidCardNumber(cardNumber),
		Network:           getNetWorkName(cardNumber),
	}

	b, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		panic(err)
	}
}

// Luhn algorithm
func isValidCardNumber(cardNumber string) bool {
	var sum int
	parity := (len(cardNumber) - 1) % 2
	digits := []rune(cardNumber)
	checkDigit, _ := strconv.Atoi(string(digits[len(digits)-1]))
	for i := len(digits) - 2; i >= 0; i-- {
		myDigit, _ := strconv.Atoi(string(digits[i]))

		if i%2 == parity {
			sum += myDigit
		} else if myDigit > 4 {
			sum += (2 * myDigit) - 9
		} else {
			sum += 2 * myDigit
		}
	}

	return checkDigit == 10-(sum%10)
}

func getRange(ranges string) (int, int) {
	tab := strings.Split(ranges, "_")
	inf, _ := strconv.Atoi(tab[0])
	sup, _ := strconv.Atoi(tab[1])

	return inf, sup
}

func getNetWorkName(cardNumber string) string {

	issuingNetwork := map[string]string{
		"34": "American Express", "37": "American Express",
		"5610": "Bankcard", "560221_560225": "Bankcard",
		"31":   "China T-Union",
		"62":   "China UnionPay",
		"36":   "Diners Club International",
		"55":   "Diners Club United States & Canada",
		"6011": "Discover Card", "644_649": "Discover Card",
		"647": "Discover Card", "622126_622925": "Discover Card",
		"60400100_60420099": "UkrCard",
		"60":                "RuPay", "81": "RuPay", "82": "RuPay", "508": "RuPay", "353": "RuPay", "356": "RuPay",
		"636":       "InterPayment",
		"637_639":   "InstaPayment",
		"3528_3589": "JCB",
		"676770":    "Maestro UK", "676774": "Maestro UK",
		"5018": "Maestro", "5020": "Maestro", "5038": "Maestro", "5893": "Maestro", "6304": "Maestro", "6761": "Maestro", "6762": "Maestro",
		"6763": "Maestro",
		"5019": "Dankort", "4571": "Dankort",
		"2200_2204": "Mir",
		"2205":      "BORICA",
		"2221_2720": "Mastercard", "51_55": "Mastercard",
		"4903": "Switch", "4905": "Switch", "4911": "Switch", "4936": "Switch", "564182": "Switch", "633110": "Switch", "6333": "Switch", "6759": "Switch | Maestro",
		"65": "Troy | Discover Card", "9792": "Troy",
		"4026": "Visa Electron", "417500": "Visa Electron", "4508": "Visa Electron", "4844": "Visa Electron", "4913": "Visa Electron", "4917": "Visa Electron",
		"1":             "UATP",
		"506099_506198": "Verve", "650002_650027": "Verve", "507865_507964": "Verve",
		"357111": "LankaPay",
		"9704":   "Napas",
	}
	for iin, net := range issuingNetwork {
		if strings.Contains(iin, "_") {
			inf, sup := getRange(iin)
			for i := inf; i <= sup; i++ {
				if strings.HasPrefix(cardNumber, strconv.Itoa(i)) {
					return net
				}
			}
		} else if strings.HasPrefix(cardNumber, iin) {
			return net
		}
	}

	if strings.HasPrefix(cardNumber, "4") {
		return "Visa"
	}

	return ""
}
