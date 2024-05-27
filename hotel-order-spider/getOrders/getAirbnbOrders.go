package getOrders

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type GetAirbnbOrderResponseBody struct {
	Data []struct {
		Listing_id                   int    `json:"listing_id"`
		Night                        int    `json:"night"`
		Confirmation_code            string `json:"confirmation_code"`
		Booked_date                  string `json:"booked_date"`
		Fullname                     string `json:"fullname"`
		Listing_name                 string `json:"listing_name"`
		Canceled_at                  string `json:"canceled_at"`
		Start_date                   string `json:"start_date"`
		End_date                     string `json:"End_date"`
		User_facing_status_localized string `json:"user_facing_status_localized"`
		Earnings                     string `json:"earnings"`
		Guest_user                   struct {
			Amount    float64 `json:"amount"`
			Currency  string  `json:"currency"`
			Full_name string  `json:"full_name"`
		} `json:"guest_user"`
	} `json:"reservations"`
	Metadata struct {
		Page_count int `json:"page_count"`
	} `json:"metadata"`
}

func GetAirbnb(platform map[string]interface{}, dateFrom, dateTo string) {
	var result string
	var resultData []ReservationsDB
	var data ReservationsDB

	cookie, ok := platform["cookie"].(string)
	if !ok {
		fmt.Println("无法获取 cookie")
	}

	// 往前一個月
	dateFromStr, _ := time.Parse("2006-01-02", dateFrom)
	oneMonthAgo := dateFromStr.AddDate(0, -1, 0)
	firstOfMonth := time.Date(oneMonthAgo.Year(), oneMonthAgo.Month(), 1, 0, 0, 0, 0, oneMonthAgo.Location())

	//第一頁
	url := `https://www.airbnb.com.tw/api/v2/reservations?locale=zh-TW&currency=MYR&_format=for_remy&_limit=40&_offset=0&collection_strategy=for_reservations_list&date_max=` + dateTo + `&date_min=` + firstOfMonth.Format("2006-01-02") + `&sort_field=start_date&sort_field=start_date&sort_order=desc&status=accepted%2Crequest%2Ccanceled`
	if err := DoRequestAndGetResponse_airbnb("GET", url, http.NoBody, cookie, &result); err != nil {
		fmt.Println("DoRequestAndGetResponse failed!")
		fmt.Println("err", err)
		return
	}
	time.Sleep(1 * time.Second)
	re := regexp.MustCompile(`(?P<currency>[^\d]+)(?P<amount>[\d,]+(\.\d+)?)`)
	var ordersData GetAirbnbOrderResponseBody
	err := json.Unmarshal([]byte(result), &ordersData)
	if err != nil {
		fmt.Println("JSON解码错误:", err)
		return
	}
	fmt.Println("頁數:", ordersData.Metadata.Page_count)

	for _, reservation := range ordersData.Data {
		if reservation.Earnings != "$0.00" && reservation.Earnings != "RM0.00" {
			data.Platform = "Airbnb"
			data.HotelId = strconv.Itoa(reservation.Listing_id)
			data.BookingId = reservation.Confirmation_code
			data.BookDate = reservation.Booked_date
			data.GuestName = reservation.Guest_user.Full_name
			data.CheckInDate = reservation.Start_date
			data.CheckOutDate = reservation.End_date

			matches := re.FindStringSubmatch(reservation.Earnings)
			if len(matches) >= 3 {
				currency := matches[1]
				amountStr := strings.ReplaceAll(matches[2], ",", "")
				amount, _ := strconv.ParseFloat(amountStr, 64)
				data.Price = amount
				if currency == "$" {
					data.Currency = "TWD"
				} else if currency == "RM" {
					data.Currency = "MYR"
				}
			}

			data.RoomNights = int64(reservation.Night)

			if reservation.User_facing_status_localized == "已確認" || reservation.User_facing_status_localized == "過往房客" || reservation.User_facing_status_localized == "為對方留下評價" || reservation.User_facing_status_localized == "等待房客評價" || reservation.User_facing_status_localized == "為對方留下評價 - 即將過期" {
				data.ReservationStatus = "已成立"
			} else if reservation.User_facing_status_localized == "由旅客取消" && data.Price != 0 {
				data.ReservationStatus = "Chargeable cancellation"
			} else if reservation.User_facing_status_localized == "由你取消" || (reservation.User_facing_status_localized == "由旅客取消" && data.Price == 0) || (reservation.User_facing_status_localized == "由 Airbnb 取消" && data.Price == 0) {
				data.ReservationStatus = "已取消"
			}

			resultData = append(resultData, data)
		}
	}
	//頁數大於一
	if ordersData.Metadata.Page_count > 1 {
		for i := 0; i < ordersData.Metadata.Page_count; i++ {
			url = `https://www.airbnb.com.tw/api/v2/reservations?locale=zh-TW&currency=MYR&_format=for_remy&_limit=40&_offset=` + strconv.Itoa((40 + (i * 40))) + `&collection_strategy=for_reservations_list&date_max=` + dateTo + `&date_min=` + firstOfMonth.Format("2006-01-02") + `&sort_field=start_date&sort_field=start_date&sort_order=desc&status=accepted%2Crequest%2Ccanceled`
			if err := DoRequestAndGetResponse_airbnb("GET", url, http.NoBody, cookie, &result); err != nil {
				fmt.Println("DoRequestAndGetResponse failed!")
				fmt.Println("err", err)
				return
			}
			time.Sleep(1 * time.Second)
			re := regexp.MustCompile(`(?P<currency>[^\d]+)(?P<amount>[\d,]+(\.\d+)?)`)
			var ordersData GetAirbnbOrderResponseBody
			err = json.Unmarshal([]byte(result), &ordersData)
			if err != nil {
				fmt.Println("JSON解码错误:", err)
				return
			}
			fmt.Println("頁數:", i+1, "/", ordersData.Metadata.Page_count)

			for _, reservation := range ordersData.Data {
				if reservation.Earnings != "$0.00" && reservation.Earnings != "RM0.00" {
					data.Platform = "Airbnb"
					data.HotelId = strconv.Itoa(reservation.Listing_id)
					data.BookingId = reservation.Confirmation_code
					data.BookDate = reservation.Booked_date
					data.GuestName = reservation.Guest_user.Full_name
					data.CheckInDate = reservation.Start_date
					data.CheckOutDate = reservation.End_date

					matches := re.FindStringSubmatch(reservation.Earnings)
					if len(matches) >= 3 {
						currency := matches[1]
						amountStr := strings.ReplaceAll(matches[2], ",", "")
						amount, _ := strconv.ParseFloat(amountStr, 64)
						data.Price = amount
						if currency == "$" {
							data.Currency = "TWD"
						} else if currency == "RM" {
							data.Currency = "MYR"
						}
					}

					data.RoomNights = int64(reservation.Night)

					if reservation.User_facing_status_localized == "已確認" || reservation.User_facing_status_localized == "過往房客" || reservation.User_facing_status_localized == "為對方留下評價" || reservation.User_facing_status_localized == "等待房客評價" || reservation.User_facing_status_localized == "為對方留下評價 - 即將過期" {
						data.ReservationStatus = "已成立"
					} else if reservation.User_facing_status_localized == "由旅客取消" && data.Price != 0 {
						data.ReservationStatus = "Chargeable cancellation"
					} else if reservation.User_facing_status_localized == "由你取消" || (reservation.User_facing_status_localized == "由旅客取消" && data.Price == 0) || (reservation.User_facing_status_localized == "由 Airbnb 取消" && data.Price == 0) {
						data.ReservationStatus = "已取消"
					}

					resultData = append(resultData, data)
				}
			}
		}
	}

	fmt.Print("resultData", resultData)

	resultDataJSON, err := json.Marshal(resultData)
	if err != nil {
		fmt.Println("JSON 轉換錯誤:", err)
		return
	}

	var resultDB string
	// 將資料存入DB
	apiurl := `http://149.28.24.90:8893/revenue_reservation/setParseHtmlToDB`
	if err := DoRequestAndGetResponse("POST", apiurl, bytes.NewBuffer(resultDataJSON), "", &resultDB); err != nil {
		fmt.Println("setParseHtmlToDB failed!")
		return
	}
}

func DoRequestAndGetResponse_airbnb(method, postUrl string, reqBody io.Reader, cookie string, resBody any) error {
	req, err := http.NewRequest(method, postUrl, reqBody)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("x-airbnb-api-key", "d306zoyjsyarp7ifhu67rjxn52tv0t20")
	req.Header.Set("Content-Type", "text/html; charset=utf-8")

	req.Header.Set("Cookie", cookie)

	client := &http.Client{Timeout: 40 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// resBody of type *string is for html
	switch resBody := resBody.(type) {
	case *string:
		// If resBody is a string
		resBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		*resBody = string(resBytes)
	default:
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, resBody); err != nil {
			return err
		}
	}

	defer resp.Body.Close()

	return nil
}
