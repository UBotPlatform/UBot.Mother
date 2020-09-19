import React from 'react';
function timeStringSince(date: Date): string {
    var r = ""
    var tick = Math.floor((new Date().valueOf() - date.valueOf()) / 1000);
    if (tick === 0){
        return "0s"
    }
    if (tick > 86400) {
        r += Math.floor(tick / 86400).toString() + "d"
        tick = tick % 86400
    }
    if (tick > 3600) {
        r += Math.floor(tick / 3600).toString() + "h"
        tick = tick % 3600
    }
    if (tick > 60) {
        r += Math.floor(tick / 60).toString() + "m"
        tick = tick % 60 
    }
    if (tick !== 0) {
        r+= tick.toString() + "s"
    }
    return r
}
interface IProps {
    since: Date
}
class TimeSince extends React.PureComponent<IProps>{
    timerID: any;
    state = {
        str: ""
    }
    componentDidMount() {
        this.tick()
        this.timerID = window.setInterval(() => this.tick(), 1000)
    }
    componentWillUnmount() {
        window.clearInterval(this.timerID);
    }
    tick() {
        this.setState({
            str: this.props.since ? timeStringSince(this.props.since) : ""
        })
    }
    render() {
        return (
            <span>{this.state.str}</span>
        )
    }
}
export default TimeSince;