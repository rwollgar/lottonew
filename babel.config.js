/*
["@babel/plugin-proposal-decorators", {"legacy": true,"loose": true}],
["@babel/plugin-proposal-class-properties", {"loose":true}],
*/
console.log('>>>> ROOT babel.config.js');
module.exports = (api) => {

    const env = api.env().toUpperCase();
    api.cache(true);

    console.log('babel.config.js: ', env);

    return {
        
        "presets": [
            ["@babel/preset-env",{"useBuiltIns":"entry","corejs":"3.18.1", "targets": {"node":"current"}}],
            "@babel/preset-react",
            "@babel/preset-typescript"
        ],
        "plugins": [
            "jsx-control-statements",
            "dynamic-import-node",
            "@babel/plugin-syntax-import-meta",
            ["@babel/plugin-proposal-decorators", {"legacy": true}],
            ["@babel/plugin-proposal-class-properties"],
            "@babel/plugin-transform-react-jsx",
            "@babel/plugin-proposal-export-namespace-from",
            "@babel/plugin-proposal-export-default-from",
            [
                "@emotion",
                {
                    // sourceMap is on by default but source maps are dead code eliminated in production
                    "sourceMap": true,
                    "autoLabel": "dev-only",
                    "labelFormat": "[local]",
                    "cssPropOptimization": true
                }
            ]
        ],
        "env": {
            "test": {
                "presets": [
                    ["@babel/preset-env",{"useBuiltIns":"entry","corejs":"3.18.1", "targets": {"node":"current"}}],
                    "@babel/preset-react",
                    "@babel/preset-typescript"
                ],
                "plugins": [
                    "@babel/plugin-transform-modules-commonjs",
                    "dynamic-import-node",
                    [
                        "@emotion",
                        {
                            // sourceMap is on by default but source maps are dead code eliminated in production
                            "sourceMap": true,
                            "autoLabel": "dev-only",
                            "labelFormat": "[local]",
                            "cssPropOptimization": true
                        }
                    ]        
                ]
            }
        }        
    }
}