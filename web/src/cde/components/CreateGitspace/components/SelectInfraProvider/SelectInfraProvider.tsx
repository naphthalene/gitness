/*
 * Copyright 2023 Harness, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import React from 'react'
import { defaultTo, isObject } from 'lodash-es'
import { Layout } from '@harnessio/uicore'
import { useFormikContext } from 'formik'
import { CDEPathParams, useGetCDEAPIParams } from 'cde/hooks/useGetCDEAPIParams'
import { OpenapiCreateGitspaceRequest, useListInfraProviderResources } from 'services/cde'
import { SelectRegion } from '../SelectRegion/SelectRegion'
import { SelectMachine } from '../SelectMachine/SelectMachine'

export const SelectInfraProvider = () => {
  const { values } = useFormikContext<OpenapiCreateGitspaceRequest>()
  const { accountIdentifier, orgIdentifier, projectIdentifier } = useGetCDEAPIParams() as CDEPathParams
  const { data } = useListInfraProviderResources({
    accountIdentifier,
    orgIdentifier,
    projectIdentifier,
    infraProviderConfigIdentifier: 'HARNESS_GCP'
  })

  const optionsList = data && isObject(data) ? data : []

  const regionOptions = optionsList
    ? optionsList.map(item => {
        return { label: defaultTo(item?.region, ''), value: defaultTo(item?.region, '') }
      })
    : []

  const machineOptions =
    optionsList
      ?.filter(item => item?.region === values?.metadata?.region)
      ?.map(item => {
        return { ...item }
      }) || []

  return (
    <Layout.Horizontal spacing="medium">
      <SelectRegion options={regionOptions} />
      <SelectMachine options={machineOptions} />
    </Layout.Horizontal>
  )
}
