const data = require('./console.json')

const dateRegex = /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z/

// convert to CSV file
const header = Object.keys(data[0]).map(v => `"${v}"`).join(',')
const csv = header + '\n' + data.map(row => Object.values(row).map(v => {
    // convert datetime from ISO to local time
    // check if it's a date
    if (dateRegex.test(v)) {
        // convert to local time
        v = new Date(v).toLocaleString()
    }
    return `"${v}"`;
}).join(',')).join('\n')

// write to file
const fs = require('fs')
fs.writeFileSync('report.csv', csv)