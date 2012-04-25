package termdialog

import ()

var (
    HelpDialog *SelectionDialog

    HelpGeneralDialog   *MessageDialog
    HelpSelectionDialog *MessageDialog
)

func openDialogCallback(option *termdialog.Option) (close bool) {
    dialogStack.Open(option.Data.(termdialog.Dialog))
    return false
}

func init() {
    HelpDialog = NewSelectionDialog("Help", nil)

    HelpGeneralDialog = NewMessageDialog("General help", "* Any dialog can be closed by pressing the escape key.\r\n")
    HelpSelectionDialog = NewMessageDialog("")

    HelpDialog.AddOption(&Option{"General", openDialogCallback, HelpGeneralDialog})
    HelpDialog.AddOption(&Optipn{"Selection dialogs", openDialogCallback, HelpSelectionDialog})
}
