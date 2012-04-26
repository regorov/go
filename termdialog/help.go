package termdialog

import ()

var (
    HelpDialog *SelectionDialog

    HelpGeneralDialog   *MessageDialog
    HelpSelectionDialog *MessageDialog
)

func openDialogCallback(option *Option) (close bool) {
    HelpDialog.GetLastDialogStack().Open(option.Data.(Dialog))
    return false
}

func init() {
    HelpDialog = NewSelectionDialog("Help", nil)

    HelpGeneralDialog = NewMessageDialog("General help", "* Any dialog can be closed by pressing the escape key.\r\n")
    HelpSelectionDialog = NewMessageDialog("Selection dialogs", "")

    HelpDialog.AddOption(&Option{"General", openDialogCallback, HelpGeneralDialog})
    HelpDialog.AddOption(&Option{"Selection dialogs", openDialogCallback, HelpSelectionDialog})
}
