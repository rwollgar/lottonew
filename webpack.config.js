const path = require('path')
//const webpack = require('webpack');
const HtmlWebpackPlugin = require('html-webpack-plugin');
//const {GitRevisionPlugin} = require('git-revision-webpack-plugin');
// const ReactRefreshWebpackPlugin = require('@pmmmwh/react-refresh-webpack-plugin');

//const {CleanWebpackPlugin} = require('clean-webpack-plugin');

//require("@babel/polyfill");
//const fs = require('fs')
//const EventsPlugin = require('webpack-event-plugin');
//const handlebars = require('handlebars');

//const env = process.env.NODE_ENV.toUpperCase();
//const gp = new GitRevisionPlugin()
const outputDir = `${path.resolve(__dirname)}/build`; //'/build';
const devPublicPath = 'http://localhost:8088/';
const srcDir = './src/ui';

//const webPath = 'web/'
//const webPathDev = 'http://localhost:9876/web/'

// console.log(`Working Directory ==> ${__dirname}`);
// console.log(`Output Directory ==> ${outputDir} Environment ==> ${env}`);

const runAfterBuildActions = (buildInfo) => {

    console.log('runAfterBuildActions: ', buildInfo.compilation.options.output.path);

    if(env === 'DEVELOPMENT') {
        console.log('\n=======AFTER Build Actions DEVELOPMENT =======\n');
    } else if (env === 'PRODUCTION') {
        console.log('\n=======AFTER Build Actions PRODUCTION =======\n');        
    }    
}

const getManifestSettings = (buildMode) => {

    console.log(`GETMANIFESTSETTINGS for Environment: ${buildMode}`);

    if(buildMode === 'PRODUCTION') {

        return {
            env: buildMode,
            version:'1.0',
            buildtime: new Date().valueOf(),
            commitHash: 'commitHash', //gp.commithash(),
            gitversion: 'gitversion', //gp.version()
        }    
    }

    if(buildMode === 'DEVELOPMENT') {

        return {
            env: buildMode,
            version:'1.0',
            buildtime: new Date().valueOf(),
            commitHash: 'commitHash (DEV)', //gp.commithash(),
            gitversion: 'gitversion (DEV)', //gp.version()
        }
    }
}

const getPlugins = (buildMode) => {

    console.log(`GETPLUGINS for Environment: ${buildMode}`);

    let plugins = [
        //new GitRevisionPlugin(),
        new HtmlWebpackPlugin({ 
            hash:true, 
            inject:'body', 
            chunks: ['appmain'],
            manifestParameters: getManifestSettings(buildMode),
            template:`${srcDir}/index_template.hbs`, 
            filename:'index.html'
        })
        //new webpack.DefinePlugin({'process.env': {'NODE_ENV': JSON.stringify(buildMode)}})

    ];

    // if(buildMode === 'DEVLOPEMNT') {
    //     plugins.push(
    //         new ReactRefreshWebpackPlugin()
    //     )
    // }
    //plugins.push(new CleanWebpackPlugin({verbose:true}));

    return plugins;
}

module.exports = (env, args) => {

    const buildMode = args.mode.toUpperCase();

    console.log(`Working Directory ==> ${__dirname} Environment ==> ${buildMode}`);
    console.log(`Output Directory ==> ${outputDir} Environment ==> ${buildMode}`);
    console.log('PublicPath ==> ', devPublicPath);

    const output = {
        hashFunction: 'xxhash64',
        filename: '[name].js',
        path: outputDir,
        //publicPath: buildMode === 'DEVELOPMENT' ? webPathDev : webPath,
        clean: true
    }

    if (buildMode === 'DEVELOPMENT') {
        output.publicPath = devPublicPath; //'http://localhost:8088/';
    }

    console.log('Output ==> ', output);
    
    return {

        target: 'web',

        performance: {hints: false},

        // optimization: {
        //     splitChunks: {
        //         cacheGroups: {
        //             default: false
        //         }
        //     }
        // },
        // entry point for app
        entry: {
            'appmain':[`${srcDir}/app.js`]
        },

        output: output,
        // output: {
        //     hashFunction: "xxhash64",
        //     filename: '[name].js',
        //     path: outputDir,
        //     publicPath: devPublicPath, //'http://localhost:8088/',
        //     //publicPath: buildMode === 'DEVELOPMENT' ? webPathDev : webPath,
        //     clean: true
        // },

        optimization: {
            minimize: buildMode === 'PRODUCTION'
        },

        module:{
            rules:
            [
                {
                    test: /\.(js|jsx)$/,
                    include: __dirname, //[path.resolve(__dirname, 'src')],
                    exclude: [path.resolve(__dirname, 'node_modules')],
                    use: {
                        loader: 'babel-loader',
                        options: {
                            plugins: [buildMode === 'DEVELOPMENT' && require.resolve('react-refresh/babel')].filter(Boolean)
                        }
                    }
                    // use:{ loader: 'swc-loader' }
                },
                {
                    test: /\.(ts|tsx)$/,
                    include: __dirname, //[path.resolve(__dirname, 'src')],
                    exclude: [path.resolve(__dirname, 'node_modules')],
                    use:{
                        loader: 'babel-loader',
                        options: {
                            plugins: [buildMode === 'DEVELOPEMNT' && require.resolve('react-refresh/babel')].filter(Boolean)
                        }
                    }
                    // use:{ loader: 'swc-loader' }
                },
                {
                    test: /\.css$/,
                    use: [ 'style-loader', 'css-loader' ]
                },
                {   test: /\.hbs$/, 
                    loader: "handlebars-loader" 
                },
                // {
                //     test: /\.(jpg|png|svg)$/,
                //     loader: 'url-loader',
                //     options: {
                //         limit: 25000,
                //         hash: "xxhash64"
                //     }
                // },
                {
                    test: /\.(jpg|png|svg)$/,
                    loader: 'file-loader',
                    options: {
                        hash: "md4",
                        limit: 25000
                    }
                }
            ],
            strictExportPresence: false

        },

        plugins: getPlugins(buildMode),
            
        resolve: {
            extensions: ['.js', '.jsx', '.ts', '.tsx']
        },

        devtool: buildMode === 'PRODUCTION' ? false : 'source-map',
        stats:'normal',
        devServer: {
            port: 8088,
            client: {
                overlay: {
                    errors: true,
                    warnings: false,
                }
            },            
            static:{
                staticOptions:{
                    contentBase: 'public'
                }
            },
            headers: {
                'Access-Control-Allow-Origin': '*'
            },            
            hot: true,
            historyApiFallback: true,
            compress: true
        }
    }
}