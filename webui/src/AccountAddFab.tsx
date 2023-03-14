import React from 'react';
import Button from '@mui/material/Button';
import TextField from '@mui/material/TextField';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import Box from '@mui/material/Box';
import Fab from '@mui/material/Fab';
import MenuItem from '@mui/material/MenuItem';
import AddIcon from '@mui/icons-material/Add';
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
                            variant="outlined"
                            margin="dense"
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
                            label="Args"
                            placeholder={"One by line, usually token or username (first line) & password (second line)"}
                            value={this.state.account_args}
                            variant="outlined"
                            margin="dense"
                            onChange={(e) => this.setState({ account_args: (e.target as HTMLInputElement).value })}
                        />
                    </DialogContent>
                    <DialogActions>
                        <Button onClick={() => this.cancel()} color="primary">Cancel</Button>
                        <Button onClick={() => this.finish()} color="primary">Finish</Button>
                    </DialogActions>
                </Dialog>
            </Box >
        );
    }
}
export default AccountAddFab