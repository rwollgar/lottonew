import React from 'react';
import m from 'moment';
import {take, takeLast} from 'ramda';

import {Typography, Button} from '@mui/material';

//import {useMatch} from 'react-location';

import styled from '@emotion/styled';

import hlp from '../../utils/helpers';
import { useLnStore } from '../store/lnstore';


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
const LnId = styled.div`
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

const LnDrawSpacer = styled.div`
    margin-left: 10px;
`

const LnDrawId = (props) => {

    return (
        <LnId>{props.children}</LnId>
    )
}

const LnDrawNumbers = (props) => {

    const { draw, std, supp } = props;
    let dummy = 0;

    return (
        
        <_lnDrawContainer id="drawcontainer">

            <a href="#" onClick={props.onDrawClick}><LnDrawId>{draw.drawid}</LnDrawId></a>
            {/* <Button variant="contained" size="small" onClick={props.onDrawClick}>{draw.drawid}</Button> */}

            {take(std,draw.numbers).map((n) => {
                if(n === 0) {
                    return (<_lnNumber key={`${draw.drawid}${++dummy}`} missing>{''}</_lnNumber>)
                } else {
                    return (<_lnNumber key={`${draw.drawid}${++dummy}`}>{n}</_lnNumber>)
                }
            })}
            <LnDrawSpacer/>            
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
    const store = useLnStore();

    const onSelectDraw = (e) => {

        e.preventDefault();
        console.log('LNDRAW: ', e.target.outerText);
        store.setSelectedDraw(e.target.outerText);

        //console.log('LNDRAW: ', id);
    }


    return (
        
        <LnDraws id="draws">
            <Typography component="div">
                {Object.values(gamedata.draws).map((d) => {    
                    return (                    
                        <LnDrawNumbers key={d.drawid}
                            draw={d}
                            std={gamedata.standardnumbers}
                            supp={gamedata.supplementary}
                            onDrawClick={onSelectDraw}/>
                        
                    )
                })}
            </Typography>
        </LnDraws>
    )
}

export default LnDraw;