package main

import (
	"bufio"
	"bytes"
	"fmt"
	"gin/getOrders"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type APIResponse struct {
	Result []RevenueAccommodationInformation `json:"result"`
}

type RevenueAccommodationInformation struct {
	HotelId                      string    `json:"hotelId" form:"hotelId" gorm:"unique;column:hotel_id;comment:旅館編號;size:512;"`
	OperationStatus              string    `json:"operationStatus" form:"operationStatus" gorm:"column:operation_status;comment:營運狀態;size:512;"`
	Country                      string    `json:"country" form:"country" gorm:"column:country;"`
	City                         string    `json:"city" form:"city" gorm:"column:city;"`
	ChineseName                  string    `json:"chineseName" form:"chineseName" gorm:"column:chinese_name;comment:旅宿中文名稱;size:512;"`
	EnglishName                  string    `json:"englishName" form:"englishName" gorm:"column:english_name;"`
	Invoice                      string    `json:"invoice" form:"invoice" gorm:"column:invoice;"`
	AccommodationType            string    `json:"accommodationType" form:"accommodationType" gorm:"column:accommodation_type;"`
	Currency                     string    `json:"currency" form:"currency" gorm:"column:currency;comment:結帳貨幣;size:512;"`
	LaunchedDate                 time.Time `json:"launchedDate" form:"launchedDate" gorm:"column:launched_date;"`
	ContactPerson                string    `json:"contactPerson" form:"contactPerson" gorm:"column:contact_person;"`
	ContactPhone                 string    `json:"contactPhone" form:"contactPhone" gorm:"column:contact_phone;"`
	ContactEmail                 string    `json:"contactEmail" form:"contactEmail" gorm:"column:contact_email;"`
	FinancialContact             string    `json:"financialContact" form:"financialContact" gorm:"column:financial_contact;"`
	FinancialPhone               string    `json:"financialPhone" form:"financialPhone" gorm:"column:financial_phone;"`
	FinancialEmail               string    `json:"financialEmail" form:"financialEmail" gorm:"column:financial_email;"`
	InvoiceEmail                 string    `json:"invoiceEmail" form:"invoiceEmail" gorm:"column:invoice_email;"`
	InvoiceTitle                 string    `json:"invoiceTitle" form:"invoiceTitle" gorm:"column:invoice_title;"`
	TaxID                        string    `json:"taxId" form:"taxId" gorm:"column:tax_id;"`
	CMS                          string    `json:"cms" form:"cms" gorm:"column:cms;"`
	CmsVendor                    string    `json:"cmsVendor" form:"cmsVendor" gorm:"column:cms_vendor;comment:串連系統(Channel manager System)的廠商;size:512;"`
	PMS                          string    `json:"pms" form:"pms" gorm:"column:pms;"`
	PMSVendor                    string    `json:"pmsVendor" form:"pmsVendor" gorm:"column:pms_vendor;"`
	SheetVersion                 string    `json:"sheetVersion" form:"sheetVersion" gorm:"column:sheet_version;"`
	SheetLink                    string    `json:"sheetLink" form:"sheetLink" gorm:"column:sheet_link;"`
	SiteminderPropertyName       string    `json:"siteminderPropertyName" form:"siteminderPropertyName" gorm:"column:siteminder_property_name;"`
	BookingAccommodationId       string    `json:"bookingAccommodationId" form:"bookingAccommodationId" gorm:"column:booking_accommodation_id;comment:Booking.com的旅宿ID;size:512;"`
	TripPrepaidAccommodationId   string    `json:"tripPrepaidAccommodationId" form:"tripPrepaidAccommodationId" gorm:"column:trip_prepaid_accommodation_id;comment:Trip.com的預付方案旅宿ID;size:512;"`
	TripPaynowAccommodationId    string    `json:"tripPaynowAccommodationId" form:"tripPaynowAccommodationId" gorm:"column:trip_paynow_accommodation_id;comment:Trip.com的現付方案旅宿ID;size:512;"`
	ExpediaAccommodationId       string    `json:"expediaAccommodationId" form:"expediaAccommodationId" gorm:"column:expedia_accommodation_id;comment:Expedia的旅宿ID;size:512;"`
	AgodaAccommodationId         string    `json:"agodaAccommodationId" form:"agodaAccommodationId" gorm:"column:agoda_accommodation_id;comment:Agoda的旅宿ID;size:512;"`
	HostelworldAccommodationId   string    `json:"hostelworldAccommodationId" form:"hostelworldAccommodationId" gorm:"column:hostelworld_accommodation_id;comment:Hostelworld的旅宿ID;size:512;"`
	Airbnb1                      string    `json:"airbnb1" form:"airbnb1" gorm:"column:airbnb1;"`
	Airbnb2                      string    `json:"airbnb2" form:"airbnb2" gorm:"column:airbnb2;"`
	Airbnb3                      string    `json:"airbnb3" form:"airbnb3" gorm:"column:airbnb3;"`
	Airbnb4                      string    `json:"airbnb4" form:"airbnb4" gorm:"column:airbnb4;"`
	Airbnb5                      string    `json:"airbnb5" form:"airbnb5" gorm:"column:airbnb5;"`
	Airbnb6                      string    `json:"airbnb6" form:"airbnb6" gorm:"column:airbnb6;"`
	Airbnb7                      string    `json:"airbnb7" form:"airbnb7" gorm:"column:airbnb7;"`
	Airbnb8                      string    `json:"airbnb8" form:"airbnb8" gorm:"column:airbnb8;"`
	Airbnb9                      string    `json:"airbnb9" form:"airbnb9" gorm:"column:airbnb9;"`
	Airbnb10                     string    `json:"airbnb10" form:"airbnb10" gorm:"column:airbnb10;"`
	Airbnb11                     string    `json:"airbnb11" form:"airbnb11" gorm:"column:airbnb11;"`
	Airbnb12                     string    `json:"airbnb12" form:"airbnb12" gorm:"column:airbnb12;"`
	Airbnb13                     string    `json:"airbnb13" form:"airbnb13" gorm:"column:airbnb13;"`
	Airbnb14                     string    `json:"airbnb14" form:"airbnb14" gorm:"column:airbnb14;"`
	Airbnb15                     string    `json:"airbnb15" form:"airbnb15" gorm:"column:airbnb15;"`
	Airbnb16                     string    `json:"airbnb16" form:"airbnb16" gorm:"column:airbnb16;"`
	Airbnb17                     string    `json:"airbnb17" form:"airbnb17" gorm:"column:airbnb17;"`
	Airbnb18                     string    `json:"airbnb18" form:"airbnb18" gorm:"column:airbnb18;"`
	Airbnb19                     string    `json:"airbnb19" form:"airbnb19" gorm:"column:airbnb19;"`
	Airbnb20                     string    `json:"airbnb20" form:"airbnb20" gorm:"column:airbnb20;"`
	Airbnb21                     string    `json:"airbnb21" form:"airbnb21" gorm:"column:airbnb21;"`
	Airbnb22                     string    `json:"airbnb22" form:"airbnb22" gorm:"column:airbnb22;"`
	Airbnb23                     string    `json:"airbnb23" form:"airbnb23" gorm:"column:airbnb23;"`
	Airbnb24                     string    `json:"airbnb24" form:"airbnb24" gorm:"column:airbnb24;"`
	Airbnb25                     string    `json:"airbnb25" form:"airbnb25" gorm:"column:airbnb25;"`
	Airbnb26                     string    `json:"airbnb26" form:"airbnb26" gorm:"column:airbnb26;"`
	Airbnb27                     string    `json:"airbnb27" form:"airbnb27" gorm:"column:airbnb27;"`
	Airbnb28                     string    `json:"airbnb28" form:"airbnb28" gorm:"column:airbnb28;"`
	Airbnb29                     string    `json:"airbnb29" form:"airbnb29" gorm:"column:airbnb29;"`
	Airbnb30                     string    `json:"airbnb30" form:"airbnb30" gorm:"column:airbnb30;"`
	Airbnb31                     string    `json:"airbnb31" form:"airbnb31" gorm:"column:airbnb31;"`
	Airbnb32                     string    `json:"airbnb32" form:"airbnb32" gorm:"column:airbnb32;"`
	Airbnb33                     string    `json:"airbnb33" form:"airbnb33" gorm:"column:airbnb33;"`
	Airbnb34                     string    `json:"airbnb34" form:"airbnb34" gorm:"column:airbnb34;"`
	Airbnb35                     string    `json:"airbnb35" form:"airbnb35" gorm:"column:airbnb35;"`
	Airbnb36                     string    `json:"airbnb36" form:"airbnb36" gorm:"column:airbnb36;"`
	Airbnb37                     string    `json:"airbnb37" form:"airbnb37" gorm:"column:airbnb37;"`
	Airbnb38                     string    `json:"airbnb38" form:"airbnb38" gorm:"column:airbnb38;"`
	Airbnb39                     string    `json:"airbnb39" form:"airbnb39" gorm:"column:airbnb39;"`
	Airbnb40                     string    `json:"airbnb40" form:"airbnb40" gorm:"column:airbnb40;"`
	Airbnb41                     string    `json:"airbnb41" form:"airbnb41" gorm:"column:airbnb41;"`
	Airbnb42                     string    `json:"airbnb42" form:"airbnb42" gorm:"column:airbnb42;"`
	Airbnb43                     string    `json:"airbnb43" form:"airbnb43" gorm:"column:airbnb43;"`
	Airbnb44                     string    `json:"airbnb44" form:"airbnb44" gorm:"column:airbnb44;"`
	Airbnb45                     string    `json:"airbnb45" form:"airbnb45" gorm:"column:airbnb45;"`
	Airbnb46                     string    `json:"airbnb46" form:"airbnb46" gorm:"column:airbnb46;"`
	Airbnb47                     string    `json:"airbnb47" form:"airbnb47" gorm:"column:airbnb47;"`
	Airbnb48                     string    `json:"airbnb48" form:"airbnb48" gorm:"column:airbnb48;"`
	Airbnb49                     string    `json:"airbnb49" form:"airbnb49" gorm:"column:airbnb49;"`
	Airbnb50                     string    `json:"airbnb50" form:"airbnb50" gorm:"column:airbnb50;"`
	OldSiteminderAccommodationId string    `json:"oldSiteminderAccommodationId" form:"oldSiteminderAccommodationId" gorm:"column:oldsite_accommodation_id;size:512;"`
	NewSiteminderAccommodationId string    `json:"newSiteminderAccommodationId" form:"newSiteminderAccommodationId" gorm:"column:newsite_accommodation_id;size:512;"`
	MastripmsAccommodationId     string    `json:"mastripmsAccommodationId" form:"mastripmsAccommodationId" gorm:"column:mastripms_accommodation_id;comment:大師的旅宿ID;size:512;"`
	OwltingAccommodationId       string    `json:"owltingAccommodationId" form:"owltingAccommodationId" gorm:"column:owlting_accommodation_id;size:512;"`
	TraiwanAccommodationId       string    `json:"traiwanAccommodationId" form:"traiwanAccommodationId" gorm:"column:traiwan_accommodation_id;size:512;"`
	CloudAccommodationId         string    `json:"cloudAccommodationId" form:"cloudAccommodationId" gorm:"column:cloud_accommodation_id;comment:雲掌櫃的旅宿ID;size:512;"`
	HostastayAccommodationId     string    `json:"hostastayAccommodationId" form:"hostastayAccommodationId " gorm:"column:hostastay_accommodation_id;size:512;"`
	BypmsAccommodationId         string    `json:"bypmsAccommodationId" form:"bypmsAccommodationId" gorm:"column:bypms_accommodation_id;comment:寶寓的旅宿ID;size:512;"`
	EndDate                      time.Time `json:"endDate" form:"endDate" gorm:"column:end_date;"`
	GroupId                      string    `json:"groupId" form:"groupId" gorm:"column:group_id;size:512;"`
	PeriodNo                     string    `json:"periodNo" form:"periodNo" gorm:"column:period_no;comment:藍新金流委託單號;size:512;"`
	MerOrderNo                   string    `json:"merOrderNo" form:"merOrderNo" gorm:"column:mer_order_no;comment:藍新金流商店訂單編號;size:512;"`
	DataSource                   string    `json:"dataSource" form:"dataSource" gorm:"column:data_source;"`
	DataDestination              string    `json:"dataDestination" form:"dataDestination" gorm:"column:data_destination;"`
	PaymentSettle                string    `json:"paymentSettle" form:"paymentSettle" gorm:"column:payment_settle;"`
	PreferredLanguage            string    `json:"preferredLanguage" form:"preferredLanguage" gorm:"column:preferred_language;"`
}

type ReservationsDB struct {
	Platform          string  `gorm:"uniqueIndex:platform_booking_id" json:"platform"`
	BookingId         string  `gorm:"uniqueIndex:platform_booking_id" json:"booking_id"`
	BookDate          string  `json:"book_date"`
	GuestName         string  `json:"guest_name"`
	NumOfGuests       int64   `json:"num_of_guests"`
	CheckInDate       string  `json:"check_in_date"`
	CheckOutDate      string  `json:"check_out_date"`
	Commission        float64 `json:"commission"`
	Price             float64 `json:"price"`
	Currency          string  `json:"currency"`
	ReservationStatus string  `json:"reservation_status"`
	NumOfRooms        int64   `json:"num_of_rooms"`
	GuestRequest      string  `json:"guest_request"`
	RoomNights        int64   `json:"room_nights"`
	HotelId           string  `json:"hotel_id"`
	Charge            string  `json:"charge"`
	RoomType          string  `json:"room_type"`
}

func main() {
	configFiles := []string{
		// "./config/config_newSIM.yaml",
		// "./config/config_oldSIM.yaml",
		// "./config/config_agoda.yaml",
		// "./config/config_expedia.yaml",
		// "./config/config_mastri.yaml",
		// "./config/config_owlting.yaml",
		// "./config/config_airbnb.yaml",
		// "./config/config_ctrip.yaml",
		"./config/config_booking.yaml",
		// "./config/config_hostelworld.yaml",
		// "./config/config_traiwan.yaml",
	}

	// 選擇執行方式
	reader, choice := choiceProcessType()

	// 設定 reqBody, startIndex, endIndex
	reqBody, startIndex, endIndex, shouldReturn1, numOfHotels := choiceFunction(choice, reader)
	if shouldReturn1 {
		return
	}

	// 取得旅館資訊
	resultDB, shouldReturn := getHotelInfoFunction(reqBody)
	if shouldReturn {
		return
	}

	for _, configFile := range configFiles {

		viper.SetConfigFile(configFile)
		viper.AddConfigPath(".")

		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Error reading config file:", err)
			return
		}

		period := viper.GetString("period")

		// period轉為時間
		timeFormat := "2006-01"
		startTime, err := time.Parse(timeFormat, period)
		if err != nil {
			fmt.Println("Error parsing period:", err)
			return
		}

		// 設定 dateFrom
		dateFrom := startTime.Format("2006-01-02")

		// 設定 dateTo
		lastDayOfMonth := startTime.AddDate(0, 1, -1)
		dateTo := lastDayOfMonth.Format("2006-01-02")

		// 執行指定 account
		shouldReturn2 := processAccount(resultDB, startIndex, endIndex, period, dateFrom, dateTo, numOfHotels)
		if shouldReturn2 {
			return
		}

	}
}

func processAccount(resultDB APIResponse, startIndex int, endIndex int, period string, dateFrom string, dateTo string, numOfHotels int) bool {
	accounts := viper.Get("account").([]interface{})
	if accounts == nil {
		fmt.Println("無法取得 account")
		return true
	}

	for _, acc := range accounts {
		account := acc.(map[string]interface{})

		accountName, ok := account["name"].(string)
		if !ok {
			fmt.Println("無法取得 account name")
			continue
		}

		fmt.Printf("accountName: %s\n", accountName)

		if platformRaw, ok := account["platform"]; ok {
			platforms, ok := platformRaw.([]interface{})
			if !ok || platforms == nil {
				fmt.Println("無法取得 platform")
				continue
			}

			for _, platformRaw := range platforms {
				platform, ok := platformRaw.(map[string]interface{})
				if !ok || platform == nil {
					fmt.Println("無法取得 platform")
					continue
				}

				platformName, ok := platform["name"].(string)
				if !ok {
					fmt.Println("無法取得 platform name")
					continue
				}

				fmt.Printf("platformName: %s\n", platformName)

				// 執行所有旅館
				processAllHotel(resultDB, startIndex, endIndex, platformName, platform, period, dateFrom, dateTo, numOfHotels)
				// 執行其他平台
				processOtherPlatform(platformName, platform, accountName, dateFrom, dateTo)

			}
		}
	}
	return false
}

func choiceProcessType() (*bufio.Reader, string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println()
	fmt.Println("--> 請選擇執行「所有旅館」或是「單一旅館」，1:所有旅館 2:單一旅館 : ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	return reader, choice
}

func processAllHotel(resultDB APIResponse, startIndex int, endIndex int, platformName string, platform map[string]interface{}, period string, dateFrom string, dateTo string, numOfHotels int) {
	for index, hotel := range resultDB.Result {
		if numOfHotels != 0 {
			if index >= startIndex && index < endIndex {
				setHotelId(index, hotel, platformName, platform, period, dateFrom, dateTo)
			}
		} else {
			setHotelId(index, hotel, platformName, platform, period, dateFrom, dateTo)
		}

	}
}

func setHotelId(index int, hotel RevenueAccommodationInformation, platformName string, platform map[string]interface{}, period string, dateFrom string, dateTo string) {
	fmt.Println()
	fmt.Println("----------------------------- START -----------------------------")
	fmt.Println("第", index+1, "間旅館")
	mrhostId := hotel.HotelId
	hotelName := hotel.ChineseName
	bookingAccommodationId := hotel.BookingAccommodationId
	agodaAccommodationId := hotel.AgodaAccommodationId
	expediaAccommodationId := hotel.ExpediaAccommodationId
	oldSIMAccommodationId := hotel.OldSiteminderAccommodationId
	newSIMAccommodationId := hotel.NewSiteminderAccommodationId
	owltingAccommodationId := hotel.OwltingAccommodationId
	operationStatus := hotel.OperationStatus

	if operationStatus != "已停運" {

		processPlatform(bookingAccommodationId, platformName, platform, period, dateFrom, dateTo, hotelName, mrhostId, agodaAccommodationId, expediaAccommodationId, oldSIMAccommodationId, newSIMAccommodationId, owltingAccommodationId)
	}
	fmt.Println()
	fmt.Println(hotelName, "執行完畢")
	fmt.Println("------------------------------ END ------------------------------")
}

func processPlatform(bookingAccommodationId string, platformName string, platform map[string]interface{}, period string, dateFrom string, dateTo string, hotelName string, mrhostId string, agodaAccommodationId string, expediaAccommodationId string, oldSIMAccommodationId string, newSIMAccommodationId string, owltingAccommodationId string) {
	if bookingAccommodationId != "" {
		if platformName == "Booking" {
			getOrders.GetBooking(platform, platformName, period, dateFrom, dateTo, bookingAccommodationId, hotelName, mrhostId)
		}
	}
	if agodaAccommodationId != "" {
		if platformName == "Agoda" {
			getOrders.GetAgoda(platform, dateFrom, dateTo, agodaAccommodationId)
		}
	}
	if expediaAccommodationId != "" {
		if platformName == "Expedia" {
			getOrders.GetExpedia(platform, dateFrom, dateTo, expediaAccommodationId)
		}
	}

	if oldSIMAccommodationId != "" {
		if platformName == "OldSIM" {
			getOrders.GetOldSIM(platform, dateFrom, dateTo, oldSIMAccommodationId)
		}
	}

	if newSIMAccommodationId != "" {
		if platformName == "NewSIM" {
			getOrders.GetNewSIM(platform, dateFrom, dateTo, newSIMAccommodationId)
		}
	}

	if owltingAccommodationId != "" {
		if platformName == "Owlting" {
			getOrders.GetOwlting(platform, dateFrom, dateTo, owltingAccommodationId)
		}
	}
}

func processOtherPlatform(platformName string, platform map[string]interface{}, accountName string, dateFrom string, dateTo string) {
	if platformName == "Ctrip" {
		getOrders.GetCtrip(platform, platformName, accountName, dateFrom, dateTo)
	}

	if platformName == "Hostelworld" {
		getOrders.GetHostelworld(platform, platformName, dateFrom, dateTo)
	}

	if platformName == "Traiwan" {
		getOrders.GetTraiwan(platform, dateFrom, dateTo)
	}

	if platformName == "MastriPMS" {
		getOrders.GetMastri(platform, dateFrom, dateTo)
	}

	if platformName == "Airbnb" {
		getOrders.GetAirbnb(platform, dateFrom, dateTo)
	}
}

func choiceFunction(choice string, reader *bufio.Reader) (string, int, int, bool, int) {
	var reqBody string
	var startIndex int
	var endIndex int
	var shouldReturn1 bool
	var numOfHotels int
	switch choice {
	case "1":
		fmt.Println("--> 請輸入每一梯次要執行幾個旅館，EX:10、20，輸入 0 則代表全部執行: ")
		numOfHotelsInput, _ := reader.ReadString('\n')
		numOfHotelsInput = strings.TrimSpace(numOfHotelsInput)
		numOfHotels, _ = strconv.Atoi(numOfHotelsInput)

		fmt.Println("每一梯次將執行", numOfHotels, "個旅館")

		if numOfHotels == 0 {
			fmt.Println("將執行所有旅館")
		} else {
			fmt.Println("--> 請輸入執行旅館梯次，EX:1、2 : ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			batchSize, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println("請重新輸入 : ", err)
				return "", 0, 0, true, numOfHotels
			}
			if batchSize == 0 {
				fmt.Println("將執行所有旅館")
			} else {
				fmt.Println("將執行第", batchSize, "梯次旅館")
			}
			startIndex = (batchSize - 1) * numOfHotels
			endIndex = batchSize * numOfHotels
		}

		reqBody = `{"mrhost_id": "","group_id":"","limit": 1000, "page_no": 1}`
	case "2":
		fmt.Println("--> 請輸入 mrhost_id，EX:R10001 : ")
		mrhostID, _ := reader.ReadString('\n')
		mrhostID = strings.TrimSpace(mrhostID)
		fmt.Println("--> 請輸入 group_id，EX:0、1 : ")
		groupID, _ := reader.ReadString('\n')
		groupID = strings.TrimSpace(groupID)
		reqBody = fmt.Sprintf(`{"mrhost_id": "%s","group_id":"%s","limit": 1000, "page_no": 1}`, mrhostID, groupID)
	default:
		fmt.Println("請重新執行 go run main.go")
		return "", 0, 0, true, numOfHotels
	}
	return reqBody, startIndex, endIndex, shouldReturn1, numOfHotels
}

func getHotelInfoFunction(reqBody string) (APIResponse, bool) {
	reqData := bytes.NewBufferString(reqBody)
	var resultDB APIResponse
	apiurl := "http://149.28.24.90:8893/rms/getHotelInfoFromDBGroup"
	if err := getOrders.DoRequestAndGetResponse("POST", apiurl, reqData, "", &resultDB); err != nil {
		fmt.Println("getHotelInfoFromDBGroup failed!", err)
		return APIResponse{}, true
	}
	return resultDB, false
}