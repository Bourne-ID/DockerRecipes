$minutesToRun = $env:RunTime
$coresToRun = $env:CoresToRun
if (!$minutesToRun) {
    $minutesToRun = 10;
}
if (!$coresToRun) {
    $coresToRun = $env:NUMBER_OF_PROCESSORS;
}

$timeend = (Get-Date) + (New-TimeSpan -m $minutesToRun)

foreach ($loopnumber in 1..$coresToRun){
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
Get-Job | Wait-Job # Wait for the launched jobs to finish