package cap

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/crypto/ssh"
	"log"
	"strings"
	"time"
)

type PasswordChecker struct {
	client       *ssh.Client
	old_password string
}

func (pc *PasswordChecker) is_pw_expired() bool {
	out, err := cleanExec(pc.client, "echo")
	if err != nil {
		log.Println("errTxt=%s", err)
	}
	log.Println("outTxt=%s", out)
	return strings.Contains(strings.ToLower(out), "expired")
}

func (pc *PasswordChecker) change_password(client ssh.Client, old_pw string, newPasswd string) error {
	log.Println("Opening shell to existing connection")
	// shell = self._ssh.invoke_shell()
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	send_pword(session, "Prompt LDAP old passwd", old_pw)
	send_pword(session, "Prompt Enter new passwd ", newPasswd)
	send_pword(session, "Prompt Reenter new passwd ", newPasswd)
	response := send_pword(session, "response", "")

	if strings.Contains(response, "updated") {
		return nil
	}
	return errors.New("Unable to update password")
}

func (pc *PasswordChecker) close() {
	// pc.ssh.close()
}

func send_pword(shell *ssh.Session, label string, password string) string {
	time.Sleep(1 * time.Second)
	// for !shell.recv_ready() {
	// 	log.Println("Waiting for: %s", label)
	// 	time.Sleep(1 * time.Second)
	// }
	// prompt := shell.recv(20000)
	// log.Println(">>>>>>>>>>>>>>>>>>>>>:  %s", prompt)
	// if password != nil {
	// 	pwd_bytes := bytes(password+"\n", "utf-8")
	// 	shell.send(pwd_bytes)
	// 	time.Sleep(1 * time.Second)
	// }
	// return prompt.decode("utf-8").strip()
	return ""
}

func NewChangePassword(
	change_cb func(new_password string)) *fyne.Container {
	old_password := widget.NewPasswordEntry()
	old_password.SetPlaceHolder("Enter old password...")
	new_password := widget.NewPasswordEntry()
	new_password.SetPlaceHolder("Enter new password...")
	// new_password.OnChanged(func() {})
	new2password := widget.NewPasswordEntry()
	new2password.SetPlaceHolder("Enter new password...")

	change := widget.NewButton("Change", func() {
		go change_cb(new_password.Text)
	})
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

        # check requirements
        # length >= 12
        # at least one lower, upper, digit, special
        # not the same as previous password
        # same as password2
*/
