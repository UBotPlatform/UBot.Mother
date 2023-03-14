import React from 'react';
import AppBar from '@mui/material/AppBar';
import Tabs from '@mui/material/Tabs';
import Tab from '@mui/material/Tab';
import Typography from '@mui/material/Typography';
import Box from '@mui/material/Box';
import Accordion from '@mui/material/Accordion';
import AccordionSummary from '@mui/material/AccordionSummary';
import AccordionDetails from '@mui/material/AccordionDetails';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import ReplayIcon from '@mui/icons-material/Replay';
import CssBaseline from '@mui/material/CssBaseline';
import Fab from '@mui/material/Fab';
import Alert from '@mui/material/Alert'
import ClientList from './ClientList';
import AccountAddFab from './AccountAddFab';
import { motherBaseURL } from './buildConfig';
import { Link } from '@mui/material';
interface TabPanelProps {
    children?: React.ReactNode;
    index: any;
    value: any;
}

function a11yProps(index: any) {
    return {
        id: `action-tab-${index}`,
        'aria-controls': `action-tabpanel-${index}`,
    };
}

function TabPanel(props: TabPanelProps) {
    const { children, value, index, ...other } = props;
    return (
        <div
            role="tabpanel"
            hidden={value !== index}
            id={`action-tabpanel-${index}`}
            aria-labelledby={`action-tab-${index}`}
            {...other}
        >
            {value === index && (
                <Box p={3}>
                    {children}
                </Box>
            )}
        </div>
    );
}

function reloadApps() {
    fetch(motherBaseURL + "/api/mother/apps/reload", {
        method: 'POST'
    })
}

export default function App() {
    const [value, setValue] = React.useState(0);

    const handleChange = (event: React.ChangeEvent<{}>, newValue: number) => {
        setValue(newValue);
    };

    return (
        <React.Fragment>
            <CssBaseline />
            <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
                <Tabs value={value} onChange={handleChange}>
                    <Tab label="Basic" {...a11yProps(0)} />
                    <Tab label="About" {...a11yProps(0)} />
                </Tabs>
            </Box>
            <Alert icon={false} severity="info">
                {"We stand against racial injustice and gender inequality that denies equal rights and opportunities."}
            </Alert>
            <TabPanel value={value} index={0}>
                <Accordion defaultExpanded variant="outlined" TransitionProps={{ unmountOnExit: true }}>
                    <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                        <Typography>Accounts</Typography>
                    </AccordionSummary>
                    <AccordionDetails>
                        <Box display="flex">
                            <ClientList endpoint="/api/mother/accounts" />
                            <AccountAddFab />
                        </Box>
                    </AccordionDetails>
                </Accordion>
                <Accordion defaultExpanded variant="outlined" TransitionProps={{ unmountOnExit: true }}>
                    <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                        <Typography>Apps</Typography>
                    </AccordionSummary>
                    <AccordionDetails>
                        <Box display="flex">
                            <ClientList endpoint="/api/mother/apps" />
                            <Fab color="primary" aria-label="reload" onClick={reloadApps}>
                                <ReplayIcon />
                            </Fab>
                        </Box>
                    </AccordionDetails>
                </Accordion>
            </TabPanel>
            <TabPanel value={value} index={1}>
                <Typography variant="h4">{"Description"}</Typography>
                <Typography variant="body1">
                    {"UBot is an open source platform for developing chat bots. It provides consistent apis for different platforms (telegram, discord and more!) in a language-independent way, making bot development WORA-able (Write once, Run anywhere)."}
                </Typography>
                <Typography variant="h4">{"Community"}</Typography>
                <Typography variant="body1">
                    {"Telegram Group: "}<Link href="https://t.me/ubotplatform">{"@ubotplatform"}</Link>{" (most recommended)"}<br />
                    {"Github Issues: "}<Link href="https://github.com/UBotPlatform">{"github.com/UBotPlatform"}</Link>
                </Typography>
                <Typography variant="h4">{"Documentation"}</Typography>
                <Typography variant="body1">
                    {"Chinese (Simplified): "}<Link href="https://www.kancloud.cn/qiqi1354092549/ubot/content">{"Read here"}</Link>
                </Typography>
            </TabPanel>
        </React.Fragment>
    );
}