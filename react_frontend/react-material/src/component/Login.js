import React, { useState } from 'react'
//import { usePostCallback } from "use-axios-react";
//import Avatar from '@material-ui/core/Avatar';
import Button from '@material-ui/core/Button';
import CssBaseline from '@material-ui/core/CssBaseline';
import TextField from '@material-ui/core/TextField';
//import FormControlLabel from '@material-ui/core/FormControlLabel';
//import Checkbox from '@material-ui/core/Checkbox';
import Link from '@material-ui/core/Link';
import Grid from '@material-ui/core/Grid';
import Box from '@material-ui/core/Box';
import { Alert }from '@material-ui/lab';
//import LockOutlinedIcon from '@material-ui/icons/LockOutlined';
import Typography from '@material-ui/core/Typography';
import {makeStyles} from '@material-ui/core/styles';
import Container from '@material-ui/core/Container';
import useAxios from 'axios-hooks'

const ROOT_URL = 'http://localhost:3000';

function refreshPage(){
    window.location.href = ROOT_URL;
}

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

export default function Login() {
    const classes = useStyles();
    const [email, setEmail] = useInput("");
    const [password, setPassword] = useInput("");

    const StatusBar = ({ data, loading, error }) => {
        console.log(data);

        if (data && data.token) {
            localStorage.setItem("token", data.token);
            refreshPage();
        }

        if (error) {
            console.log(error);
        }

        return (
            <span>
                {loading && <Alert severity="info">Loading...</Alert>}
                {error && <Alert severity="info">Error sending requests...</Alert>}
            </span>
        )
    };

    const [{ data, loading, error }, executeTokenPost] = useAxios({
            url: "http://localhost:8080/api/token",
            method: 'post',
        },
        { manual: true }
    );

    function createToken() {
        executeTokenPost({
            data: {
                email: email,
                password: password,
            }
        });
    }

    return (
        <Container component="main" maxWidth="xs">
            <CssBaseline/>
            <div className={classes.paper}>

                <StatusBar className={classes.margin} loading={loading} error={error} data={data} />

                <Typography component="h1" variant="h5" className={classes.margin}>
                    Login
                </Typography>

                    <Grid container spacing={2}>
                        <Grid item xs={12}>
                            <TextField
                                variant="outlined"
                                required
                                fullWidth
                                id="email"
                                label="Email Address"
                                name="email"
                                autoComplete="email"
                                value={email}
                                onChange={setEmail}
                            />
                        </Grid>
                        <Grid item xs={12}>
                            <TextField
                                variant="outlined"
                                required
                                fullWidth
                                name="password"
                                label="Password"
                                type="password"
                                id="password"
                                autoComplete="current-password"
                                value={password}
                                onChange={setPassword}
                            />
                        </Grid>
                    </Grid>

                <Button
                    fullWidth
                    variant="contained"
                    color="primary"
                    className={classes.submit}
                    onClick={() => {createToken()}}
                >
                    Login Now
                </Button>
                    <Grid container justify="flex-end">
                        <Grid item>
                            <Link href="/signup" variant="body2">
                                Do not have an account? Sign up Now
                            </Link>
                        </Grid>
                    </Grid>

            </div>
            <Box mt={5}>
            </Box>
        </Container>
    );
}