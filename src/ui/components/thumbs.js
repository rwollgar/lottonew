import React, {useState, useRef} from 'react';

import {IconButton} from '@mui/material';
import ThumbUpIcon from '@mui/icons-material/ThumbUpRounded';
import ThumbDownIcon from '@mui/icons-material/ThumbDownRounded';

import styled from '@emotion/styled';

const ScoSurveyPage = styled.div`
    height: calc(100% - 58px);
    width: 100%;
    background-color: white; 
    border-radius: 0 0 8px 8px;
    display: flex;
    flex-direction: column;
    justify-content: center;
    box-sizing: border-box;
    padding: 20px;
`

const ScoSurveyTagline = styled.div`
    font-family:inherit;
    font-size: ${props => {return props.fontsize ? props.fontsize : 20;}}px;
    font-weight: ${props => {return props.fontweight ? props.fontweight : 100;}};
    text-align: center;
    margin-top: ${props => {return props.margintop ? props.margintop : 0;}}px;
`

const ScoThumbUpDown = styled.div`
    display: flex;
    flex-direction: row;
    justify-content: center;
    background-color: transparent;
    margin: 20px 0 50px 0;
    height: 70px;
`

const ScoThumbUp = styled(IconButton)`

    color: ${props => props.thumbstate === 'up' ? '#38DD98' : 'grey'};
    background: ${props => props.thumbstate === 'up' ? 'rgba(56, 221, 152, 0.1)' : 'transparent'};
    border: ${props => props.thumbstate === 'up' ? '2px solid #38DD98' : '2px solid transparent'};
    height: 5.5rem;
    border-radius: 50%!important;

`

const ScoThumbDown = styled(IconButton)`

    color: ${props => props.thumbstate === 'down' ? '#FF335C' : 'grey'};
    background: ${props => props.thumbstate === 'down' ? 'rgba(255, 51, 92, 0.1)' : 'transparent'};
    border: ${props => props.thumbstate === 'down' ? '2px solid #FF335C' : '2px solid transparent'};
    height: 5.5rem;
    border-radius: 50%!important;

`

const ThumbsUpDown = (props) => {

    const [thumbState, setThump] = useState('');

    let disabled = true;    

    const surveyResultUp = () => {
        //console.log('surveyResultGood');
        setThump('up');
        disabled = false;
    }

    const surveyResultDown = () => {
        //console.log('surveyResultBad');
        setThump('down');
        disabled = false;
    }


    return (
        <ScoSurveyPage>
            <ScoSurveyTagline fontsize={18} margintop={30}>
                How did we do?
            </ScoSurveyTagline>
            <ScoThumbUpDown>
                <ScoThumbUp disableRipple disableFocusRipple thumbstate={thumbState} onClick={surveyResultUp}>
                    <ThumbUpIcon style={{width:'4rem', height:'4rem'}}/>
                </ScoThumbUp>
                <ScoThumbDown disableRipple disableFocusRipple thumbstate={thumbState} onClick={surveyResultDown}>
                    <ThumbDownIcon color="green" style={{width:'4rem', height:'4rem'}}/>
                </ScoThumbDown>
            </ScoThumbUpDown>
        </ScoSurveyPage>

    )
}

export default ThumbsUpDown