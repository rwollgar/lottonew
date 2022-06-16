import React, {useEffect} from "react";
//import ReactDOM from 'react-dom';
import { createRoot } from 'react-dom/client';
import { Outlet, ReactLocation, Router } from "@tanstack/react-location";
//import { ReactLocationDevtools } from 'react-location-devtools';
import useAsyncEffect from 'use-async-effect';

import { Box, Paper, Typography } from '@mui/material';
import styled from '@emotion/styled';

//import { ThemeProvider, createTheme } from '@mui/material/styles';
//import styled from 'styled-components';

//import LnGamePanel from "./components/gamepanel";
//import LnDraw from './components/lndraw';
import LnLayout from './components/layout';
import { useLnStore, dataApi } from './store/lnstore';
//import lnStore from './store/lnStore';

import HomePage from './pages/home';
import GamePage from './pages/game';

// const OuterContainer = styled(Box)`
//     display: flex;
//     flex-direction: column;
//     width: 80%;
//     height: ${props => props.height};
//     margin-top: 0;
//     margin-bottom: 0;
//     padding:1 2 1 2;
//     border-style: solid;
//     border-width: 2px;
//     margin-left: auto;
//     margin-right: auto;
//     min-height: 650px;
// `
// const Titlebar = styled(Box)`
//     border-style: solid;
//     border-width: 2px;
//     border-color: blue;
//     padding: 2 4 2 4;
// `

// const Commandbar = styled(Box)`
//     border-style: solid;
//     border-width: 2px;
//     margin-top:2px;
//     margin-bottom: 2px;
//     border-color: green;
//     padding: 2 4 2 4;
// `

// const Footer = styled(Box)`
//     border-style: solid;
//     border-width: 2px;
//     margin-top: 2px;
//     margin-bottom: 2px;
//     border-color: seagreen;
//     padding: 2 4 2 4;
// `

// const Content = styled(Paper)`
//     flex: 1;
//     display: flex;
//     flex-direction: column;
//     padding: .5rem;
//     border-style: solid;
//     border-width: 2px;
//     overflow-y: auto;
//     border-color: purple;
// `

// const LnContainer = styled(Box)`

//     display: flex;
//     flex-wrap: wrap;
//     justify-content: center;
//     flex-direction: row;
//     padding-top: 10px;
//     padding-bottom: 10px;
//     width: 95%;

// `

// const LnGames = (props) => {

//     const store = useLnStore();
//     const imageMap = store.imageMap;

//     return (
//         <LnContainer>
//             {store.data.map((g) => {
//                 return (
//                     <LnGamePanel key={g.name} game={g} image={imageMap[g.name]}/>
//                 )
//             })}
//         </LnContainer>
//     )
// }

const Bill = styled(Box)`
    display:flex;
    flex-direction: column;
    width: 300px;
    height: 300px;
    border: black solid 1px;
    font-family: Roboto,Helvetica,Arial,sans-serif;
`

const BillHeader = styled(Box)`
    height: 10%;
    width: 100%;
    background-color: blue;
`

const BillBody = styled(Box)`
    display: flex;
    flex-direction:row;
    height: 90%;
    width: 100%;
`

const BillBodyLabel = styled(Box)`
    display: flex;
    flex-direction: row;
    width: 40%;
    height: 100%;
    background-color: yellow;
`
const BillBodyValue = styled(Box)`
    display: flex;
    flex-direction: row;
    width: 60%;
    height: 100%;
    background-color: purple;
`

const location = new ReactLocation();

const routes = [
    {
        path: '/',
        element: <HomePage/>
    },
    {
        path: 'gamedetails',
        children: [{
            path: ':gameid',
            element: <GamePage />,
            loader: async ({ params: { gameid }, parentMatch }) => (
                {
                    gamedata: await dataApi.getGameData(gameid)
                }
            )
        }]
    }
]


const App = (props) => {

    const store = useLnStore();

    // useAsyncEffect(async () => {
    //     console.log('App: Hello from useEffect');
    //     await store.initStore();
    // }, []);

    //console.log(props)

    return (
        <Router routes={routes} location={location}>
            {/* <OuterContainer id="outercontainer"height="98%">
                <Titlebar id="titlebar"><Typography>TITLE</Typography></Titlebar>
                <Commandbar id="commandbar"><Typography>COMMANDS</Typography></Commandbar>

                <Content id="content"> */}
                <LnLayout id="layout" error={store.lastError}>
                    <Outlet />

                    <Bill>
                        <BillHeader>HEADER</BillHeader>
                        <BillBody>
                            <BillBodyLabel>LABEL</BillBodyLabel>
                            <BillBodyValue>VALUE</BillBodyValue>
                        </BillBody>
                    </Bill>

                </LnLayout>
                

                {/* </Content>
                <Footer id="footer"><Typography>FOOTER</Typography></Footer>
            </OuterContainer> */}
        </Router>
    )
}

// ReactDOM.render(<App />, document.getElementById('root'), () => {
//     console.log('APP rendered')
// });

//const container = document.getElementById('root');
const AppBoot = (props) => {

    const store = useLnStore();

    useAsyncEffect(async () => {
        await store.initStore();
        console.log('APPBOOT: Lotto App rendered.');
    },[]);

    return (<App/>)
}

const root = createRoot(document.getElementById('root'));
root.render(<AppBoot/>);

if (module.hot) {
    module.hot.accept();
    console.log("MODULE.HOT");
}
// if (module.hot) {
//     module.hot.accept('./app', () => {
//         const NextApp = require('./app').default; // Get the updated code
//         render(NextApp);
//     });
// }