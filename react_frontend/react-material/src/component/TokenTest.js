import React, { useState, useEffect } from 'react'
import Button from '@material-ui/core/Button';
import CssBaseline from '@material-ui/core/CssBaseline';
import TextField from '@material-ui/core/TextField';
import Link from '@material-ui/core/Link';
import Grid from '@material-ui/core/Grid';
import Box from '@material-ui/core/Box';
import { Alert }from '@material-ui/lab';
import Typography from '@material-ui/core/Typography';
import {makeStyles} from '@material-ui/core/styles';
import Container from '@material-ui/core/Container';
import useAxios from 'axios-hooks'

function useInput(initialValue) {
    const [value, setValue] = useState(initialValue);

    function handleChange(e){
        setValue(e.target.value);
    }

    return [value, handleChange];
}

const useStyles = makeStyles(theme => ({
    paper: {
        marginTop: theme.spacing(8),
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
    },
    avatar: {
        margin: theme.spacing(1),
        backgroundColor: theme.palette.secondary.main,
    },
    form: {
        width: '100%', // Fix IE 11 issue.
        marginTop: theme.spacing(3),
    },
    submit: {
        margin: theme.spacing(3, 0, 2),
    },
    margin: {
        marginTop: theme.spacing(3),
        marginBottom: theme.spacing(3),
    },
}));

export default function TokenTest() {
    const classes = useStyles();
    const [token, setToken] = useState("");

    useEffect(() => {
        let token = localStorage.getItem("token");
        setToken(token);
    }, [])

    const StatusBar = ({ data, loading, error }) => {
        if (error) {
            console.log(error);
        }
        return (
            <span>
                {!error && data && <Alert severity="info">Request sent successfully</Alert>}
                {loading && <Alert severity="info">Loading and waiting for response...</Alert>}
                {error && <Alert severity="info">Error sending requests...</Alert>}
            </span>
        )
    };

    const [{ data, loading, error }, executeTokenHello] = useAxios({
            url: "http://localhost:8080/api/hello",
            method: 'get',
        },
        { manual: true }
    );

    function getHello(token) {
        executeTokenHello({
            headers: {'Authorization': 'Bearer ' + token}
        });
    }

    return (
        <Container component="main" maxWidth="xs">
            <CssBaseline/>
            <div className={classes.paper}>

                <StatusBar className={classes.margin} loading={loading} error={error} data={data} />

                <Grid container spacing={2}>
                </Grid>

                <Button
                    fullWidth
                    variant="contained"
                    color="primary"
                    className={classes.submit}
                    onClick={() => {getHello(token)}}
                >
                    Hello With Valid Token
                </Button>

                <Button
                    fullWidth
                    variant="contained"
                    color="primary"
                    className={classes.submit}
                    onClick={() => {getHello("")}}
                >
                    Hello With Error Token
                </Button>

            </div>
        </Container>
    );
}