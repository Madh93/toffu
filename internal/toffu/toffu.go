package toffu

import (
	"errors"
	"fmt"
	"time"

	"github.com/Madh93/toffu/internal/woffuapi"
	"github.com/spf13/viper"
)

type Toffu struct {
	api   *woffuapi.WoffuAPI
	signs *Signs
}

func New() *Toffu {
	return &Toffu{}
}

func NewWithAPI(api *woffuapi.WoffuAPI) *Toffu {
	return &Toffu{api: api}
}

func (t *Toffu) SetAPI(api *woffuapi.WoffuAPI) {
	t.api = api
}

func (t *Toffu) SetSigns(signs woffuapi.Signs) {
	t.signs = &Signs{signs: signs}
}

func (t *Toffu) GenerateToken() (err error) {
	if viper.GetString("woffu_token") != "" {
		timestamp, erro := getTokenExpiration()
		if erro != nil {
			return erro
		}
		expiration := time.Unix(timestamp, 0).Format("Monday, 02 January 2006 15:04:05 MST")
		fmt.Printf("A Woffu token is already valid until %s! Skipping...\n", expiration)
		return
	}

	fmt.Println("Generating a new Woffu token...")

	if t.api == nil {
		username, password := getCredentials()
		t.api = woffuapi.NewWithBasicAuth(username, password)
	}

	token, err := t.api.CreateToken()
	if err != nil {
		return
	}

	viper.Set("woffu_token", token.AccessToken)
	if err = viper.WriteConfig(); err != nil {
		return
	}

	fmt.Printf("\nA Woffu token has been generated!\n")
	return
}

func (t *Toffu) ClockIn() (err error) {
	fmt.Println("Trying to clock in...")

	err = t.initAPIWithToken()
	if err != nil {
		return err
	}

	userId, err := getUserId()
	if err != nil {
		return
	}

	signs, err := t.api.GetSigns()
	if err != nil {
		return err
	}
	t.SetSigns(signs)

	if t.signs.hasAlreadyClockedIn() {
		return errors.New("error clocking in, you have already clocked in")
	}

	workday, err := t.api.GetUserWorkDay(userId)
	if err != nil {
		return err
	}

	if workday.ScheduleHours <= 0.0 {
		return errors.New("error clocking in, no scheduled working hours today")
	}

	err = t.api.Sign(userId)
	if err != nil {
		return fmt.Errorf("error clocking in: %v", err)
	}

	fmt.Println("You have clocked in sucessfully!")

	return
}

func (t *Toffu) ClockOut() (err error) {
	fmt.Println("Trying to clock out...")

	err = t.initAPIWithToken()
	if err != nil {
		return err
	}

	userId, err := getUserId()
	if err != nil {
		return
	}

	signs, err := t.api.GetSigns()
	if err != nil {
		return err
	}
	t.SetSigns(signs)

	if !t.signs.hasAlreadyClockedIn() {
		return errors.New("error clocking out, you have not clocked in or you have already clocked out")
	}

	err = t.api.Sign(userId)
	if err != nil {
		return fmt.Errorf("error clocking out: %v", err)
	}

	fmt.Println("You have clocked out sucessfully!")

	return
}

func (t *Toffu) GetStatus() (err error) {
	err = t.initAPIWithToken()
	if err != nil {
		return err
	}

	slots, err := t.api.GetSignSlots()
	if err != nil {
		return err
	}

	// Current status
	status := ""

	if len(slots) > 0 && (woffuapi.SignSlotEvent{}) == slots[len(slots)-1].Out {
		status = "In Office"
	} else {
		status = "Out of Office"
	}

	fmt.Printf("Status: %s\n", status)

	// Hours worked
	totalDuration := 0 * time.Second
	location, err := time.LoadLocation("Europe/Madrid") // TODO: Woffu TZ or Office TZ?
	if err != nil {
		return err
	}

	for _, slot := range slots {
		inTime, _ := time.Parse("15:04:05", slot.In.ShortTrueTime)
		outTime, _ := time.Parse("15:04:05", time.Now().In(location).Format("15:04:05"))
		// In Office
		if slot.Out.ShortTrueTime != "" {
			outTime, _ = time.Parse("15:04:05", slot.Out.ShortTrueTime)
		}
		// Day/Night shift transition
		if inTime.After(outTime) {
			outTime = outTime.Add(24 * time.Hour)
		}
		delta := outTime.Sub(inTime)
		totalDuration += delta
	}

	// Remaining hours
	userId, err := getUserId()
	if err != nil {
		return
	}

	workday, err := t.api.GetUserWorkDay(userId)
	if err != nil {
		return err
	}

	remainingDuration := time.Duration(workday.ScheduleHours*float64(time.Hour)) - totalDuration

	// Show hours worked and remaining hours
	fmt.Printf("Total hours worked today: %s", secondsToHumanReadable(totalDuration))

	if remainingDuration > 0 {
		fmt.Printf(" (%s remaining)\n", secondsToHumanReadable(remainingDuration))
	} else {
		fmt.Println("")
	}

	return
}

func (t *Toffu) initAPIWithToken() (err error) {
	if t.api == nil {
		token, err := getToken()
		if err != nil {
			return err
		}
		t.api = woffuapi.NewWithToken(token)
	}

	return
}
