import React from 'react';
import './App.css';
import DrawerNavbar from "./component/DrawerNavbar";
import Container from '@material-ui/core/Container';
import {BrowserRouter as Router, Route, Switch} from "react-router-dom";
import ListUserComponent from "./component/user/ListUserComponent";
import AddUserComponent from "./component/user/AddUserComponent";
import EditUserComponent from "./component/user/EditUserComponent";
import SignUp from "./component/SignUp";
import Login from "./component/Login";
import TokenTest from "./component/TokenTest";
import Typography from "@material-ui/core/Typography";
import Link from "@material-ui/core/Link";

function App() {

    const Copyright = () => {
        return (
            <Typography variant="body2" color="textSecondary" align="center">
                {'Copyright © '}
                <Link color="inherit" href="https://material-ui.com/">
                    Your Website
                </Link>{' '}
                {new Date().getFullYear()}
                {'.'}
            </Typography>
        );
    }

    return (
        <div>
            <Container>
                <Router>
                    <DrawerNavbar/>
                    <Switch>
                        <Route path="/" exact component={ListUserComponent} />
                        <Route path="/users" component={ListUserComponent} />
                        <Route path="/add-user" component={AddUserComponent} />
                        <Route path="/edit-user" component={EditUserComponent} />
                        <Route path="/tokentest" component={TokenTest} />
                        <Route path="/signup" component={SignUp} />
                        <Route path="/login" component={Login} />
                    </Switch>
                </Router>
                {Copyright()}
            </Container>
        </div>
    );
}

export default App;
