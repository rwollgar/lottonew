import React from 'react';
import m from 'dayjs';
import styled from '@emotion/styled';

import {Box, CardHeader, CardMedia, Divider, Typography} from '@mui/material';
import {Link} from '@tanstack/react-location';

import hlp from '../../utils/helpers';

// self.$RefreshReg$ = () => {};
// self.$RefreshSig$ = () => () => {};

//import { Games } from '@material-ui/icons';


// const MyPaper = styled(Paper)`

//     width:350px;
//     height:250px;
//     background-color:lightyellow;
//     flex: 0 0 350px; 
//     margin-right: 5px;
//     margin-bottom: 5px;

// `
// const BallImage = styled.img`
//     position:absolute;
//     clip:rect(0px,350px,100px,0px);
// `
// .css-bhp9pd-MuiPaper-root-MuiCard-root {
// 	background-color: #fff;
// 	color: rgba(0, 0, 0, 0.87);
// 	-webkit-transition: box-shadow 300ms cubic-bezier(0.4, 0, 0.2, 1) 0ms;
// 	transition: box-shadow 300ms cubic-bezier(0.4, 0, 0.2, 1) 0ms;
// 	border-radius: 4px;
// 	box-shadow: 0px 2px 1px -1px rgba(0,0,0,0.2),0px 1px 1px 0px rgba(0,0,0,0.14),0px 1px 3px 0px rgba(0,0,0,0.12);
// 	overflow: hidden;
// }

const LnCard = styled(Box)`
    background: ${hlp.env.theme.gamepanelbg};
    margin: 10px;
    width: 450px;
    color: black;
    text-transform: capitalize;
    border-radius: 4px;
    box-shadow: 0px 2px 1px -1px rgba(0,0,0,0.2),0px 1px 1px 0px rgba(0,0,0,0.14),0px 1px 3px 0px rgba(0,0,0,0.12);
    overflow: hidden;
`

const LnHeader = styled.div`
    display: flex;
    flex-direction: row;
    justify-content: space-between;
`

const LnHeaderItem = styled.div`
    font-family: "Roboto","Helvetica","Arial",sans-serif;
    display: flex;
    flex-direction: row;
    justify-content: ${props => props.justify ? props.justify : 'flex-start'};
    height: 30px;
`

const LnHeaderSubItem = styled.div`
    border-style: solid;
    border-color: ${hlp.env.theme.bordercolor};
    border-width: 0.1rem;
    font-size: 1.4rem;
    padding-left: 2px;
    padding-right: 2px;
    margin-left: 3px;
    font: inherit;
`

const LnCardMedia = styled(CardMedia)` 
    width: 450px;
    height: 103px;
    margin-bottom: 10px;
`
const LnContent = styled.div`
    display: flex;
    flex-direction: column;
    padding-left: 16px;
    padding-right: 30px;
`

const LnContentItem = styled.div`
    display: flex;
    font-weight: 500;
    flex-direction: row;
    justify-content: space-between;
    background-color: transparent;
`

const LnLabel = styled.div`
    width: 40%;
    text-transform: none;
`
const LnValue = styled.div`
    width: 60%;
`

const LnDivider = styled(Divider)`
    width: 100%;
    margin-left: 20%;
    height: 2px;
`
// const LnCardActions = styled(CardActions)`

//     justify-content: center;

// `

const Header = (props) => {

    const {game} = props;

    return (

        <LnHeader>
            <LnHeaderItem>
                {/* <a href={`/gamedetails/${game.name}`}>{game.name.replace('-',' ')}</a> */}
                <Link to={`/gamedetails/${game.name}`}>{game.name.replace('-',' ')}</Link>
            </LnHeaderItem>
            <LnHeaderItem justify="flex-end">
                <LnHeaderSubItem>{game.standardnumbers}</LnHeaderSubItem>
                <LnHeaderSubItem>{game.supplementary}</LnHeaderSubItem>
                <LnHeaderSubItem>{game.maxnumber}</LnHeaderSubItem>
                <LnHeaderSubItem>{game.maxsupplementary}</LnHeaderSubItem>
            </LnHeaderItem>
        </LnHeader>
    )
}

const LnGamePanel = (props) => {

    const {game, image} = props;

    console.log(game);

    return (

        <Typography component="div">
            <LnCard>
                <CardHeader title={<Header game={game}/>}/>
                <LnCardMedia component="img" src={image} title={game.name}/>
                <LnContent>
                    {/* <LnDivider/> */}
                    <LnContentItem>
                        <LnLabel>Last draw</LnLabel>
                        <LnValue><b>{`${m(game.lastdraw.date).format('ddd, D MMM YYYY')}`}</b></LnValue>
                    </LnContentItem>
                    <LnContentItem>
                        <LnLabel>Next draw</LnLabel>
                        <LnValue><b>{`${m(game.lastdraw.date).add(7, 'day').format('ddd, D MMM YYYY')}`}</b></LnValue>
                    </LnContentItem>
                    <LnContentItem>
                        <LnLabel>Draws</LnLabel>
                        <LnValue><b>{game.draws ? game.Draw.length: game.numdraws}</b></LnValue>
                    </LnContentItem>
                    <LnDivider/>
                </LnContent>
            </LnCard>
        </Typography>
    )
}

export default LnGamePanel;
