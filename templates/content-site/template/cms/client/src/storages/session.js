import createStorage from 'sham-ui-data-storage';

export const { storage, useStorage } = createStorage( {
    name: '',
    email: '',
    sessionValidated: false,
    isAuthenticated: false,
    isSuperuser: false
}, {
    DI: 'session:storage'
} );
