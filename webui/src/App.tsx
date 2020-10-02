import React from 'react';
import { makeStyles, Theme } from '@material-ui/core/styles';
import AppBar from '@material-ui/core/AppBar';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';
import Typography from '@material-ui/core/Typography';
import Box from '@material-ui/core/Box';
import Accordion from '@material-ui/core/Accordion';
import AccordionSummary from '@material-ui/core/AccordionSummary';
import AccordionDetails from '@material-ui/core/AccordionDetails';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import ReplayIcon from '@material-ui/icons/Replay';
import Container from '@material-ui/core/Container';
import Fab from '@material-ui/core/Fab';
import Alert from '@material-ui/lab/Alert'
import ClientList from './ClientList';
import AccountAddFab from './AccountAddFab';
import { motherBaseURL } from './buildConfig';
import { Link } from '@material-ui/core';
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

const useStyles = makeStyles((theme: Theme) => ({
    root: {
        flexGrow: 1,
        backgroundColor: theme.palette.background.paper,
    },
}));

function reloadApps() {
    fetch(motherBaseURL + "/api/mother/apps/reload", {
        method: 'POST'
    })
}

export default function App() {
    const classes = useStyles();
    const [value, setValue] = React.useState(0);

    const handleChange = (event: React.ChangeEvent<{}>, newValue: number) => {
        setValue(newValue);
    };

    return (
        <Container className={classes.root}>
            <AppBar position="static">
                <Tabs value={value} onChange={handleChange}>
                    <Tab label="Basic" {...a11yProps(0)} />
                    <Tab label="About" {...a11yProps(0)} />
                </Tabs>
            </AppBar>
            <Alert icon={false} severity="info">
                {"We stand against racial injustice and gender inequality that denies equal rights and opportunities."}
            </Alert>
            <TabPanel value={value} index={0}>
                <Accordion defaultExpanded variant="outlined" TransitionProps={{ unmountOnExit: true }}>
                    <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                        <Typography>Accounts</Typography>
                    </AccordionSummary>
                    <AccordionDetails>
                        <ClientList endpoint="/api/mother/accounts" />
                        <AccountAddFab />
                    </AccordionDetails>
                </Accordion>
                <Accordion defaultExpanded variant="outlined" TransitionProps={{ unmountOnExit: true }}>
                    <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                        <Typography>Apps</Typography>
                    </AccordionSummary>
                    <AccordionDetails>
                        <ClientList endpoint="/api/mother/apps" />
                        <Fab color="primary" aria-label="reload" onClick={reloadApps}>
                            <ReplayIcon />
                        </Fab>
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
        </Container>
    );
}