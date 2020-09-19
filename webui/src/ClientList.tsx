import React from 'react';
import Typography from '@material-ui/core/Typography';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import CardActions from '@material-ui/core/CardActions';
import Button from '@material-ui/core/Button';
import Box from '@material-ui/core/Box';
import Grid from '@material-ui/core/Grid';
import Chip from '@material-ui/core/Chip';
import { ServiceStatus, ServiceInfo, startService, stopService } from './serviceOperation';
import { motherBaseURL } from './buildConfig'
import TimeSince from './TimeSince'
interface ClientInfo {
    id: string
    binded_service?: ServiceInfo
}
interface IState {
    clients: ClientInfo[]
}
interface IProps {
    endpoint: string
}
class ClientList extends React.Component<IProps, IState> {
    timerID: any;
    state = {
        clients: []
    } as IState
    componentDidMount() {
        this.tick()
    }
    componentWillUnmount() {
        window.clearTimeout(this.timerID);
    }
    async tick() {
        try {
            const res = await fetch(motherBaseURL + this.props.endpoint);
            const json = await res.json();
            if (json) {
                this.setState({
                    clients: json
                });
            } else {
                this.setState({
                    clients: []
                });
            }
        } catch (error) {

        }
        this.timerID = window.setTimeout(
            () => this.tick(),
            1000
        );
    }
    render() {
        return (
            <Grid container spacing={3}>
                {
                    this.state.clients?.map((client) =>
                        <Grid item key={client.id}>
                            <Card raised>
                                <CardContent>
                                    <Typography variant="h5" component="h2">
                                        {client.id}
                                    </Typography>
                                    {client.binded_service &&
                                        <Typography variant="body2" color="textSecondary">Service: {client.binded_service.id}</Typography>
                                    }
                                    {client.binded_service?.status === ServiceStatus.Running &&
                                        <Typography variant="body2">
                                            {"Uptime: "}
                                            <TimeSince since={new Date(client.binded_service.launch_at)} />
                                        </Typography>
                                    }
                                </CardContent>
                                <CardActions>
                                    <Grid justify="space-between" alignItems="center" container>
                                        <Grid item>
                                            {client.binded_service?.status === ServiceStatus.Stopped &&
                                                <Chip size="small" label="Stopped" color="secondary" />
                                            }
                                            {client.binded_service?.status === ServiceStatus.Running &&
                                                <Chip size="small" label="Running" color="primary" />
                                            }
                                            {client.binded_service?.status === ServiceStatus.Starting &&
                                                <Chip size="small" label="Starting" color="primary" />
                                            }
                                            {client.binded_service?.status === ServiceStatus.Exited &&
                                                <Chip size="small" label="Exited" color="secondary" />
                                            }
                                        </Grid>
                                        <Grid item>
                                            {client.binded_service ?
                                                <Box>
                                                    <Button
                                                        size="small"
                                                        color="primary"
                                                        onClick={() => client.binded_service && startService(client.binded_service)}>Start</Button>
                                                    <Button
                                                        size="small"
                                                        color="primary"
                                                        onClick={() => client.binded_service && stopService(client.binded_service)}>Stop</Button>
                                                </Box> :
                                                <Typography>Self maintained</Typography>}
                                        </Grid>
                                    </Grid>
                                </CardActions>
                            </Card>
                        </Grid>
                    )
                }
            </Grid>
        );
    }
}
export default ClientList;