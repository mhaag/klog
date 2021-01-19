package mac_widget

import (
	"fmt"
	"klog/app"
	"klog/lib/caseymrm/menuet"
	"klog/record"
	"klog/service"
	"os/exec"
	"time"
)

var blinker = 1

func render(ctx app.Context) []menuet.MenuItem {
	var items []menuet.MenuItem

	items = append(items, func() []menuet.MenuItem {
		if rs, file, err := ctx.Bookmark(); err == nil {
			return renderRecords(rs, file)
		}
		return []menuet.MenuItem{{
			Text:       "No input file specified",
			FontWeight: menuet.WeightBold,
		}, {
			Text: "Specify one by running:",
		}, {
			Text: "klog widget --file yourfile.klg",
		}}
	}()...)

	items = append(items, menuet.MenuItem{
		Type: menuet.Separator,
	}, menuet.MenuItem{
		Text:  "Start widget on login",
		State: hasLaunchAgent(ctx.HomeDir()),
		Clicked: func() {
			var err error
			if hasLaunchAgent(ctx.HomeDir()) {
				err = removeLaunchAgent(ctx.HomeDir())
			} else {
				err = createLaunchAgent(ctx.HomeDir())
			}
			if err != nil {
				fmt.Println(err)
			}
		},
	})

	return items
}

func renderRecords(records []record.Record, file app.File) []menuet.MenuItem {
	var items []menuet.MenuItem
	now := time.Now()
	nowTime, _ := record.NewTime(now.Hour(), now.Minute())
	nowDate := record.NewDateFromTime(now)
	today := func() record.Record {
		candidates, _ := service.FindFilter(records, service.Filter{
			BeforeEq: nowDate, AfterEq: nowDate,
		})
		if len(candidates) == 1 {
			return candidates[0]
		}
		return nil
	}()

	totalTimeValue, isRunningCurrently := func() (string, bool) {
		if today == nil {
			return "", false
		}
		if today.OpenRange() != nil {
			result := service.HypotheticalTotal(today, nowTime).ToString()
			result += "  "
			switch blinker {
			case 1:
				result += "◷"
			case 2:
				result += "◶"
			case 3:
				result += "◵"
			case 4:
				result += "◴"
				blinker = 0
			}
			blinker++
			return result, true
		} else {
			return service.Total(today).ToString(), false
		}
	}()

	items = append(items, menuet.MenuItem{
		Text: file.Name,
		Children: func() []menuet.MenuItem {
			total := service.Total(records...)
			should := service.ShouldTotal(records...)
			diff := record.NewDuration(0, 0).Minus(should).Plus(total)
			plus := ""
			if diff.InMinutes() > 0 {
				plus = "+"
			}
			items := []menuet.MenuItem{
				{
					Text: "Show in Finder...",
					Clicked: func() {
						cmd := exec.Command("open", file.Location)
						_ = cmd.Start()
					},
				},
				{Type: menuet.Separator},
				{Text: "Total: " + total.ToString()},
				{Text: "Should: " + should.ToString() + "!"},
				{Text: "Diff: " + plus + diff.ToString()},
			}
			showMax := 7
			for i, r := range service.Sort(records, false) {
				if i == 0 {
					items = append(items, menuet.MenuItem{Type: menuet.Separator})
				}
				if i == showMax {
					items = append(items, menuet.MenuItem{Text: fmt.Sprintf("(%d more)", len(records)-showMax)})
					break
				}
				items = append(items, menuet.MenuItem{Text: r.Date().ToString() + ": " + service.Total(r).ToString()})
			}
			return items
		},
	})

	items = append(items, menuet.MenuItem{
		Text:  "Track time",
		State: isRunningCurrently,
		Clicked: func() {
			if isRunningCurrently {
				//ctx.QuickStopAt(nowDate, nowTime)
			} else {
				//ctx.QuickStartAt(nowDate, nowTime)
			}
		},
	})

	if today != nil {
		items = append(items, menuet.MenuItem{
			Text: "Today: " + totalTimeValue,
		})
	}

	return items
}