package telegram

const msgHelp = `I can save and keep you pages. Just send me URL.
To get all the pages from your collection, use the /all command.
In order to get random page from your collection, use /rnd command.
Caution! After that, this URL will be removed from your collection.`

const msgHello = "Hi there! \n\n" + msgHelp

const (
	msgUnknownCommand = "Unknown command 🧐"
	msgNoSavedPages   = "You have no saved pages 😕"
	msgSaved          = "Saved! 👌"
	msgAlreadyExists  = "You have already this page in your list 😌"
)
