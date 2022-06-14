export function fileSize( bytes ) {
    const thresh = 1000;
    if ( isNaN( bytes ) || bytes === undefined ) {
        bytes = 0;
    }
    if ( Math.abs( bytes ) < thresh ) {
        return bytes + ' B';
    }
    const units = [
        'kB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'
    ];
    let u = -1;
    do {
        bytes /= thresh;
        ++u;
    } while ( Math.abs( bytes ) >= thresh && u < units.length - 1 );
    return bytes.toFixed( 1 ) + ' ' + units[ u ];
}

export function duration( totalSeconds ) {
    const days = Math.floor( totalSeconds / 86400 );
    totalSeconds %= 86400;
    const hours = Math.floor( totalSeconds / 3600 );
    totalSeconds %= 3600;
    const minutes = Math.floor( totalSeconds / 60 );
    const seconds = totalSeconds % 60;
    return `${days} days ${hours} hours ${minutes} minutes ${seconds} seconds`;
}
