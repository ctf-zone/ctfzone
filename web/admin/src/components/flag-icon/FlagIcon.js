import * as React from 'react'
import FlagIconFactory, { functions } from 'react-flag-icon-css'

const FlagIcon = FlagIconFactory(React, { useCssModules: false })

export const countries = functions.countries.getCountries()
export default FlagIcon
