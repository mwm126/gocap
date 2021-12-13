package login

import (
	"errors"
	"log"
	"strings"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	// "aeolustec.com/capclient/cap"
)

func NewChangePassword(change_cb func(new_password string)) *fyne.Container {
	old_password := widget.NewPasswordEntry()
	old_password.SetPlaceHolder("Enter old password...")
	new_password := widget.NewPasswordEntry()
	new_password.SetPlaceHolder("Enter new password...")
	new2password := widget.NewPasswordEntry()
	new2password.SetPlaceHolder("Enter new password...")

	change := widget.NewButton("Change", func() {
		go change_cb(new_password.Text)
	})

	check_new_password := func(_ string) {
		result := password_passes(old_password.Text, new_password.Text, new2password.Text)
		if result == nil {
			change.Enable()
		} else {
			log.Println(result)
			change.Disable()
		}
	}
	new_password.OnChanged = check_new_password
	new2password.OnChanged = check_new_password
	change.Disable()
	return container.NewVBox(old_password, new_password, new2password, change)
}

/*

class ChangePasswordDialog(QDialog):
    def __init__(self, parent: QWidget, old_passwd: str, prompt: str) -> None:
        QDialog.__init__(self, parent)

        # variables
        self.requirements = [False] * 7
        self.old_passwd = old_passwd
        self.success = False

        # load the ui
        self.ui = ui = Ui_Form()
        self.ui.setupUi(self)

        ui.lineEdit_newpass.textChanged.connect(self.newpass_changed)
        ui.lineEdit_confirm_pass.textChanged.connect(self.confirm_changed)

        ui.pushButton_change_pass.clicked.connect(self.handle_change)
        ui.pushButton_cancel.clicked.connect(self.handle_cancel)

        ui.label_prompt.setText(prompt)

        self.req_text_list = [
            ui.label_char,
            ui.label_lower,
            ui.label_upper,
            ui.label_num,
            ui.label_sym,
        ]

        # hide widgets
        ui.label_not_prev.setVisible(False)
        ui.label_not_match.setVisible(False)

        ui.lineEdit_newpass.setFocus()
*/

func password_passes(old, new, new2 string) error {
	if new != new2 {
		return errors.New("Passwords do not match")
	}
	if old == new {
		return errors.New("Password is the same as the previous password")
	}
	if len(new) < 12 {
		return errors.New("Password must have length >=12 characters")
	}
	if !strings.ContainsAny(new, "abcdefghijklmnopqrstuvwxyz") {
		return errors.New("Password must contain a lowercase letter")
	}
	if !strings.ContainsAny(new, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return errors.New("Password must contain an uppercase letter")
	}
	if !strings.ContainsAny(new, "0123456789") {
		return errors.New("Password must contain a digit")
	}
	return nil
}
