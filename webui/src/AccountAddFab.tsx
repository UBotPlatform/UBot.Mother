import React from 'react';
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import Box from '@material-ui/core/Box';
import Fab from '@material-ui/core/Fab';
import MenuItem from '@material-ui/core/MenuItem';
import AddIcon from '@material-ui/icons/Add';
import { motherBaseURL } from './buildConfig';
interface AccountProviderInfo {
    id: string
}
class AccountAddFab extends React.Component {
    state = {
        open: false,
        account_type: "",
        account_args: "",
        account_providers: [] as Array<AccountProviderInfo>
    }
    async open() {
        try {
            const res = await fetch(motherBaseURL + "/api/mother/account_providers");
            const json = await res.json();
            if (json != null) {
                this.setState({
                    open: true,
                    account_type: "",
                    account_args: "",
                    account_providers: json
                });
            }
        } catch (error) {

        }
    };
    cancel() {
        this.setState({
            open: false
        });
    }
    finish() {
        var args = this.state.account_args.match(/[^\r\n]+/g); 
        fetch(motherBaseURL + "/api/mother/accounts", {
            body: JSON.stringify([
                {
                    "type": this.state.account_type, 
                    "args": args
                }
            ]),
            headers: {
                'content-type': 'application/json'
            },
            method: 'POST'
        })
        this.setState({
            open: false
        });
    }
    render() {
        return (
            <Box>
                <Fab color="primary" aria-label="add" onClick={() => this.open()}>
                    <AddIcon />
                </Fab>
                <Dialog open={this.state.open} aria-labelledby="account-add-dialog-title">
                    <DialogTitle id="account-add-dialog-title">Add Account</DialogTitle>
                    <DialogContent>
                        <DialogContentText>
                            {"Type the account type and args (usually token or username & password)"}
                        </DialogContentText>
                        <TextField
                            autoFocus
                            fullWidth
                            select
                            id="account_type"
                            label="Type"
                            value={this.state.account_type}
                            onChange={(e) => this.setState({ account_type: (e.target as HTMLInputElement).value })}
                        >
                            {this.state.account_providers.map((x) => (
                                <MenuItem key={x.id} value={x.id}>{x.id}</MenuItem>
                            ))}
                        </TextField>
                        <TextField
                            autoFocus
                            fullWidth
                            multiline
                            rows={4}
                            id="account_args"
                            label="Args (one per line)"
                            value={this.state.account_args}
                            onChange={(e) => this.setState({ account_args: (e.target as HTMLInputElement).value })}
                        />
                    </DialogContent>
                    <DialogActions>
                        <Button onClick={() => this.cancel()} color="primary">Cancel</Button>
                        <Button onClick={() => this.finish()} color="primary">Finish</Button>
                    </DialogActions>
                </Dialog>
            </Box>
        );
    }
}
export default AccountAddFab