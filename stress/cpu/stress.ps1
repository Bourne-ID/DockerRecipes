$minutesToRun = $env:RunTime
if (!$minutesToRun) {
    $minutesToRun = 1;
}

$timeend = (Get-Date) + (New-TimeSpan -m $minutesToRun)

foreach ($loopnumber in 1..$env:NUMBER_OF_PROCESSORS){
    Start-Job -ScriptBlock{ param($timeend)
    $result = 1
    foreach ($number in 1..2147483647){
        $result = $result * $number
        if ((Get-Date) -gt $timeend) {
            break;
        }
    }# end foreach
    } -Arg $timeend # end Start-Job
}# end foreach
Get-Job | Wait-Job