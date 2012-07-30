package main

import (
    "github.com/mattn/go-gtk/gdk"
//    "github.com/mattn/go-gtk/glib"
    "github.com/mattn/go-gtk/gtk"
    "os"
    "sync"
//    "time"
)

var (
    exitChannelsLock sync.Mutex
    exitChannels = make([]chan bool, 0)
)

var instructionSetDialog *gtk.GtkWindow

func main() {
    gdk.ThreadsInit()
    gdk.ThreadsEnter()
    gtk.Init(nil)
    
    createGUI()
    
    go gtkExitListener()
    
    /* // Testing globalExit
    go func() {
        time.Sleep(time.Second * 10)
        globalExit()
    }()
    */
    
    gtk.Main()
    for {}
}

func globalExit() {
    exitChannelsLock.Lock()
    for _, ch := range exitChannels {
        ch <- true
    }
    exitChannelsLock.Unlock()
    
    exitChannelsLock.Lock()
    for _, ch := range exitChannels {
        <-ch // Expect a reply
    }
    exitChannelsLock.Unlock()
    
    os.Exit(0)
}

func gtkExitListener() {
    exitChannelsLock.Lock()
    exitChan := make(chan bool)
    exitChannels = append(exitChannels, exitChan)
    exitChannelsLock.Unlock()
    
    <-exitChan
    gtk.MainQuit()
    //time.Sleep(time.Second * 2)
    exitChan <- true
}

func createGUI() {
    window := gtk.Window(gtk.GTK_WINDOW_TOPLEVEL)
    window.SetPosition(gtk.GTK_WIN_POS_CENTER)
    window.SetTitle("CPU Lab")
    window.SetSizeRequest(800, 600)
    window.Connect("destroy", globalExit)
    
    mainVBox := gtk.VBox(false, 1)
    window.Add(mainVBox)
    
    menubar := gtk.MenuBar()
    mainVBox.PackStart(menubar, false, false, 0)
    
    fileMenuItem := gtk.MenuItemWithMnemonic("_File")
    menubar.Append(fileMenuItem)
    
    fileMenu := gtk.Menu()
    fileMenuItem.SetSubmenu(fileMenu)
    
    quitMenuItem := gtk.MenuItemWithMnemonic("_Quit")
    quitMenuItem.Connect("activate", globalExit)
    fileMenu.Append(quitMenuItem)
    
    cpuMenuItem := gtk.MenuItemWithMnemonic("_CPU")
    menubar.Append(cpuMenuItem)
    
    cpuMenu := gtk.Menu()
    cpuMenuItem.SetSubmenu(cpuMenu)
    
    instructionSetMenuItem := gtk.MenuItemWithMnemonic("_Instruction Set...")
    instructionSetMenuItem.Connect("activate", func() {instructionSetDialog.ShowAll()})
    cpuMenu.Append(instructionSetMenuItem)
    
    window.ShowAll()
    
    // ---------------------------------------------------------------------------------------------
    
    instructionSetDialog = gtk.Window(gtk.GTK_WINDOW_TOPLEVEL)
    instructionSetDialog.SetPosition(gtk.GTK_WIN_POS_CENTER)
    instructionSetDialog.SetTitle("Instruction Set")
    instructionSetDialog.SetSizeRequest(400, 300)
    instructionSetDialog.Connect("delete-event", func() (ret bool) {instructionSetDialog.HideAll(); return true})
    
    instructionSetVBox := gtk.VBox(false, 1)
    instructionSetDialog.Add(instructionSetVBox)
    
    instructionSetScroller := gtk.ScrolledWindow(nil, nil)
    instructionSetVBox.PackStart(instructionSetScroller, true, true, 0)
    
    instructionSetStore := gtk.ListStore(gtk.GTK_TYPE_STRING, gtk.GTK_TYPE_STRING) // Name, format assigned
    
    instructionSetView := gtk.TreeView()
    instructionSetView.SetModel(instructionSetStore)
    instructionSetView.AppendColumn(gtk.TreeViewColumnWithAttributes("Name", gtk.CellRendererText(), "text", 0))
    instructionSetView.AppendColumn(gtk.TreeViewColumnWithAttributes("Format assigned", gtk.CellRendererText(), "text", 1))
    instructionSetScroller.Add(instructionSetView)
    
    instructionSetButtonBox := gtk.HBox(false, 1)
    instructionSetVBox.PackStart(instructionSetButtonBox, false, false, 0)
    
    instructionSetAddButton := gtk.ButtonWithLabel("Add")
    instructionSetAddButton.Clicked(func() {
        var iter gtk.GtkTreeIter
        instructionSetStore.Append(&iter)
        instructionSetStore.Set(&iter, "<unnamed>", "no")
    })
    instructionSetButtonBox.PackStart(instructionSetAddButton, false, false, 0)
    
    instructionSetRemoveButton := gtk.ButtonWithLabel("Remove")
    instructionSetRemoveButton.Clicked(func() {
        var iter *gtk.GtkTreeIter
        var path *gtk.GtkTreePath
        var column *gtk.GtkTreeViewColumn
        instructionSetView.GetCursor(&path, &column)
        instructionSetStore.GetIter(iter, path)
        instructionSetStore.Remove(iter)
    })
    instructionSetButtonBox.PackStart(instructionSetRemoveButton, false, false, 0)
}
