import React from "react";

import { Box, Paper, Typography } from '@mui/material';


import styled from '@emotion/styled';

//import styled from 'styled-components';


const Titlebar = styled(Box)`
    border-style: solid;
    border-width: 2px;
    border-color: blue;
    padding: 2 4 2 4;
`

const Commandbar = styled(Box)`
    border-style: solid;
    border-width: 2px;
    margin-top:2px;
    margin-bottom: 2px;
    border-color: green;
    padding: 2 4 2 4;
`

const Footer = styled(Box)`
    border-style: solid;
    border-width: 2px;
    margin-top: 2px;
    margin-bottom: 2px;
    border-color: seagreen;
    padding: 2 4 2 4;
    height: 25px;
`

const Content = styled(Paper)`
    flex: 1;
    display: flex;
    flex-direction: column;
    padding: .5rem;
    border-style: solid;
    border-width: 2px;
    overflow-y: auto;
    border-color: purple;
`
const OuterContainer = styled(Box)`
    display: flex;
    flex-direction: column;
    width: 80%;
    height: ${props => props.height};
    margin-top: 0;
    margin-bottom: 0;
    padding:1 2 1 2;
    border-style: solid;
    border-width: 2px;
    margin-left: auto;
    margin-right: auto;
    min-height: 650px;
`

const LnLayout = (props) => {

    //console.log("LNLAYOUT:", props);

    return (

        <OuterContainer id="outercontainer" height="98%">
            <Titlebar id="titlebar"><Typography>TITLE</Typography></Titlebar>
            <Commandbar id="commandbar"><Typography>COMMANDS</Typography></Commandbar>

            <Content id="content">
                {props.children}    
            </Content>

            <Footer id="footer"><Typography>{props.error}</Typography></Footer>
        </OuterContainer>

    )
}

export default LnLayout;