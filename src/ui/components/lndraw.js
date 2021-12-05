import React from 'react';
import m from 'moment';
import {take, takeLast} from 'ramda';

import {Typography} from '@mui/material';

import {useMatch} from 'react-location';

import styled from '@emotion/styled';

import hlp from '../../utils/helpers';

 //(Container)
const LnDraws = styled.div`
    background-color: inherit;
    max-width: 460px;
    overflow-x : hidden;
    overflow-y: scroll;
    border: solid lightgrey 2px;
    border-radius: 5px;
`
//(Container)
const _lnDrawContainer = styled.div`
    display: flex;
    flex-direction: row;
    border-radius: 5px;
    &:hover {
        background: ${hlp.env.theme.gamepanelbg};
    }    
`
const _lnNumber = styled.div`
    margin: 3px;
    width: 32px;
    height: 32px;
    border: ${props => props.missing ? '0px' : '1px'} solid black;
    border-radius: 4px;
    padding: 4px 2px 1px 2px;
    font-size: 20px;
    font-weight: ${props => props.supp ? 600 : 400};
    text-align: center;
    color: ${props => props.supp ? 'green' : 'black'};
    &:hover {
        background: ${hlp.env.theme.numberbg};
    }    
`
const _lnDrawid = styled.div`
    margin: 3px 10px 3px 3px;
    width: 60px;
    height: 32px;
    border: 1px solid black;
    border-radius: 4px;
    padding: 1px 2px 1px 2px;
    font-size: 20px;
    font-weight: 500;
    text-align: center;
`

const _spacer = styled.div`
    margin-left: 10px;
`

const LnDrawId = (props) => {

    return (
        <_lnDrawid>{props.children}</_lnDrawid>
    )
}

const LnDrawNumbers = (props) => {

    const {draw, std, supp} = props;
    var dummy = 99;

    const linkClick = (e) => {
        e.preventDefault();
        console.log(e.target.outerText);
    }

    return (
        
        <_lnDrawContainer id="drawcontainer">

            <a href="#" onClick={linkClick}><LnDrawId>{draw.drawid}</LnDrawId></a>

            {take(std,draw.numbers).map((n) => {
                if(n === 0) {
                    return (<_lnNumber key={`${draw.drawid}${++dummy}`} missing>{''}</_lnNumber>)
                } else {
                    return (<_lnNumber key={`${draw.drawid}${n}`}>{n}</_lnNumber>)
                }
            })}
            <_spacer/>            
            {takeLast(supp,draw.numbers).map((n) => {
                return (<_lnNumber supp key={`${draw.drawid}${n}`}>{n}</_lnNumber>)
            })}
        </_lnDrawContainer>
    )

}

const LnDraw = (props) =>  {

    //console.log('LnDraw: ', props.details);
    //console.log('LnDraw: ', props.details.draws.length);
    const { gamedata } = props.data;

    return (
        
        <LnDraws id="draws">
            <Typography component="div">
                {Object.values(gamedata.draws).map((d) => {    
                    return (                    
                        <LnDrawNumbers  key={d.drawid} 
                                        draw={d} 
                                        std={gamedata.standardnumbers} 
                                        supp={gamedata.supplementary}/>
                    )
                })}
            </Typography>
        </LnDraws>
    )
}

export default LnDraw;