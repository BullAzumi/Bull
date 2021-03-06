/* Created by Bull and Notriv
 * I wrote this version and it was a lot of work! So please keep that in mind when criticizing.
 * Notriv has made the version clearer and made it bug free as far as possible. 
 * A big thanks to Notriv! Buddy without you, the script would probably never have come this far!
 */
 
/*
 * VERSION 1.2
 */
 

/* DESCRIPTION
 * Automatically bid on auctions
 * Maximum resources adjustable
 * Setting only available resources
 * Registered players will not be outbid
 */

/*---------------------------------------------------------------------------------------------------------------------------*/



//######################################## SETTINGS START ########################################

highestMet = 3000000                                    //what is the maximum bid of metal
highestCrys = 1500000                                   //what is the maximum bid of crystal
highestDeut = 1000000                                   //what is the maximum bid of deuterium
bidHome = "M:4:363:4"                                   //from which planet should be bid?
playerIgnore = ["Imperator Bla", "Bull", "Notriv"]      //add player names that should not be outbid

//######################################## SETTINGS END ########################################


ownPlayerID = GetCachedPlayer().PlayerID
tmpMet = 0
tmpCrys = 0
tmpDeut = 0
highestBid = 0
totalItems = 0
celt = GetCachedCelestial(bidHome)
if celt == nil {
    LogError(bidHome + " is not one of your planet/moon")
    return
}

if highestMet == 0 && highestCrys == 0 && highestDeut == 0{
    LogError("Highest bid not Set!")
       return
}

func AucDo(ress) {

    bid = { celt.GetID() : ress }
    return DoAuction(bid)
}

func ressDefine(mustBid) {

    switch mustBid {
        case tmpMet > 0 && mustBid <= tmpMet:
            LogDebug("Use metal")
            tmpRess = NewResources(mustBid, 0, 0)
            tmpMet -= mustBid
            return tmpRess
        case tmpMet > 0 && mustBid > tmpMet && mustBid - tmpMet <= tmpCrys * 1.5:
            LogDebug("Use metal and crystal")
            tmpRess = NewResources(tmpMet, Round((mustBid - tmpMet) / 1.5) + 1, 0)
            tmpCrys -= Round((mustBid - tmpMet) / 1.5) + 1
            tmpMet -= tmpMet
            return tmpRess
        case tmpMet > 0 && mustBid > tmpMet && mustBid - tmpMet - (tmpCrys * 1.5) <= tmpDeut * 3:
            LogDebug("Use metal and crystal and deuterium")
            tmpRess = NewResources(tmpMet, tmpCrys, Round((mustBid - tmpMet - (tmpCrys * 1.5)) / 3) + 1)
            tmpDeut -= Round((mustBid - tmpMet - (tmpCrys * 1.5)) / 3) + 1
            tmpCrys -= tmpCrys
            tmpMet -= tmpMet 
            return tmpRess
        case tmpMet > 0 && mustBid > tmpMet && mustBid - tmpMet <= tmpDeut * 3:
            LogDebug("Use metal and deuterium")
            tmpRess = NewResources(tmpMet, 0, Round((mustBid - tmpMet) / 3) + 1)
            tmpDeut -= Round((mustBid - tmpMet) / 3) + 1
            tmpMet -= tmpMet
            return tmpRess
        case tmpCrys > 0 && mustBid / 1.5 <= tmpCrys:
            LogDebug("Use crystal")
            tmpRess = NewResources(0, Round(mustBid / 1.5) + 1, 0)
            tmpCrys -= Round(mustBid / 1.5) + 1
            return tmpRess
        case tmpCrys > 0 && mustBid / 1.5 > tmpCrys && mustBid - (tmpCrys * 1.5) <= tmpDeut:
            LogDebug("Use crystal and deuterium")
            tmpRess = NewResources(0, tmpCrys, Round((mustBid - (tmpCrys * 1.5)) / 3) + 1)
            tmpDeut -= Round((mustBid - (tmpCrys * 1.5)) / 3) + 1
            tmpCrys -= tmpCrys
            return tmpRess
        case tmpDeut > 0 && mustBid / 3 <= tmpDeut:
            LogDebug("Use deuterium")
            tmpRess = NewResources(0, 0, Round(mustBid / 3) + 1)
            tmpDeut -= Round(mustBid / 3) + 1
            return tmpRess
        default:
            LogWarn("The bid is to high! No Resources left!")
            return nil
    } 
}

func refreshTime(TimeEnd) {
    switch TimeEnd {     
        case TimeEnd <= 300:                    //5 min
        LogDebug("Only 5 min")
        return Random(2, 5)

        case TimeEnd <= 600:                    //10 min
        LogDebug("Only 10 Min")                        
        return Random(60, 120)

        case TimeEnd <= 1800:                   //30 min
        LogDebug("Only 30 Min")                        
        return Random(180, 300)

        case TimeEnd <= 3600:                   //60 min
        LogDebug("Only 60 Min")                        
        return Random(300, 600)

        default:
        LogError("Unknown TimeEnd value", TimeEnd)
        return Random(5, 10)
    }
}

func customSleep(sleepTime) {
    if sleepTime <= 0 {
        sleepTime = Random(5, 10)
    }
    LogInfo("Wait " + ShortDur(sleepTime))
    Sleep(sleepTime * 1000)
}

func didWon(auc) {
    if auc.HighestBidderUserID == ownPlayerID {
        LogInfo("You won the auction with " + Dotify(auc.CurrentBid) + " resources!")
        LogInfo(auc.CurrentItem + " has been added to your inventory!")
        totalItems++
        LogInfo("You won total " + totalItems)
    }
}

func resetTmpRess(){
    ress, _ = celt.GetResources()
    if ress.Metal < highestMet && ress.Crystal < highestCrys && ress.Deuterium < highestDeut {
        LogError("No resources available at " + bidHome + "! Change bidHome and restart the script!")
        StopScript(__FILE__)
    }
    if ress.Metal >= highestMet{
        tmpMet = highestMet
        LogDebug(Dotify(tmpMet) + " metal has been set")
    }else {
        LogWarn("You have not enough metal on " + bidHome)
    }
    if ress.Crystal >= highestCrys{
        tmpCrys = highestCrys
        LogDebug(Dotify(tmpCrys) + " crystal has been set")
    }else {
        LogWarn("You have not enough crystal on " + bidHome)
    }
    if ress.Deuterium >= highestDeut{
        tmpDeut = highestDeut
        LogDebug(Dotify(tmpDeut) + " deuterium has been set")
    }else {
        LogWarn("You have not enough deuterium on " + bidHome)
    }

    highestBid = tmpMet + Round(tmpCrys * 1.5) + (tmpDeut * 3)
    LogDebug("Resources has been reset!")
    LogDebug("Maximum bid set do " + Dotify(highestBid))
}

func processAuction() {
    auc, err = GetAuction()
    if err != nil {
        LogError(err)
        return Random(5, 10)
    }
    if auc.HasFinished {
        if auc.Endtime > 3600 {
            LogInfo("There is currently no auction")
        } else {
            LogInfo("Auction has finished")
        }
        didWon(auc)
        resetTmpRess()
        return auc.Endtime + 10
    }
    for player in playerIgnore {
        if auc.HighestBidder == player {
            LogInfo("This is a friend we cannot beat!")
            return refreshTime(auc.Endtime)
        }
    }
    if auc.HighestBidderUserID == ownPlayerID {
        LogInfo("You are the highest bidder!")
        return refreshTime(auc.Endtime)
    }
    if auc.MinimumBid > highestBid {
        LogInfo("Resources exceeded! Wait until the next auction!")
        return auc.Endtime + 10
    }

    ress = ressDefine(auc.MinimumBid - auc.AlreadyBid)
    LogInfo("You are not the highest bidder! Bid " + ress + " resources!")
    err = AucDo(ress)
    if err != nil {
        LogError(err)
        return Random(5, 10)
    }
    return refreshTime(auc.Endtime)
}

func doWork() {
    resetTmpRess()
    for { // forever process auctions
        sleepTime = processAuction()
        customSleep(sleepTime)
    }
}
doWork()
