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
import { Redirect } from "react-router-dom";

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

/*
const ConvertedText = (data) => {
    console.log(data);
    return (
        <Box>test</Box>
    )
}
 */

export default function SignUp() {
    const classes = useStyles();

    const [email, setEmail] = useInput("");
    const [password, setPassword] = useInput("");
    //const [queryText, setQueryText] = useInput("");
    //const [queryText, setQueryText] = useInput("");
    //const query = {id: "", querytext: queryText}
    //const query = {email:email, password:password};

    /*
    function postRequest({ email, password }) {
        var url = "http://127.0.0.1:8080/api/user";
        return {
            url: url,
            data: { email, password}
        };
    }
     */

    /*
    function getRequest() {
        var url = "http://127.0.0.1:8080/api/hello";
        return {
            url: url,
        }
    }
     */

    const StatusBar = ({ data, loading, error }) => {
        var signupOK = false;
        if (data) {
            signupOK = true;
        }

        if (error) {
            console.log(error);
        }

        return (
            <span>
                {signupOK ? <Redirect to="/login" /> : null}
                {loading && <Alert severity="info">Loading...</Alert>}
                {error && <Alert severity="info">Error sending requests...</Alert>}
            </span>
        )
    };

    const [{ data, loading, error }, executePost] = useAxios({
            url: "http://localhost:8080/api/user",
            method: 'post',
        },
        { manual: true }
    );

    function createUser() {
        executePost({
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
                    Sign up
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
                    onClick={() => {createUser()}}
                >
                    Create My Account
                </Button>
                    <Grid container justify="flex-end">
                        <Grid item>
                            <Link href="/login" variant="body2">
                                Already have an account? Login in
                            </Link>
                        </Grid>
                    </Grid>

            </div>
            <Box mt={5}>
            </Box>
        </Container>
    );
}