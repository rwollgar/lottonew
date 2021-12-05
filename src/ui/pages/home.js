import React from "react";

import { Box, Paper, Typography } from '@mui/material';
import styled from '@emotion/styled';

import LnGamePanel from "../components/gamepanel";
import { useLnStore } from '../store/lnstore';

const LnContainer = styled(Box)`

    display: flex;
    flex-wrap: wrap;
    justify-content: center;
    align-content: flex-start;
    flex-direction: row;
    padding-top: 10px;
    padding-bottom: 10px;
    width: 95%;

`

const HomePage = (props) =>{

    const store = useLnStore();
    const imageMap = store.imageMap;

    return (
        <LnContainer>
            {store.data.map((g) => {    
                return (
                    <LnGamePanel key={g.name} game={g} image={imageMap[g.name]}/>
                )
            })}
            
        </LnContainer>
    )
}

export default HomePage;