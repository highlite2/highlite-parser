package sylius

// CreateTaxon creates a taxon
/*func (c *client) CreateTaxon(ctx context.Context, taxon *transfer.TaxonNew) (*transfer.Taxon, error) {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	url := c.getURL("/v1/taxons/")

	request, err := c.getRequest(ctx)
	if err != nil {
		return nil, err
	}

	result := &transfer.TaxonRaw{}
	resp, err := request.SetResult(result).Get(url)

	if err != nil {
		return nil, fmt.Errorf("request to %s failed: %s", url, err.Error())
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("request to %s ended with status %s", url, resp.Status())
	}

	return transfer.ConvertRawTaxon(result), nil
}

func (c *client) getRequest(ctx context.Context) (*resty.Request, error){
	token, err := c.getToken()
	if err != nil {
		return nil, err
	}

	request := resty.R().SetContext(ctx).SetHeader("Authorization", "Bearer "+token)

	return request, nil
}

func (c *client) postRequest(ctx context.Context, url string, data interface{}) error {
	request, err := c.getRequest(ctx)
	if err != nil {
		return err
	}

	request.SetResult(data).SetC
}*/
