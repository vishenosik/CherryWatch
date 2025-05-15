package api

func IsInfo(code int) bool { return code >= 100 && code <= 199 }

func IsSuccess(code int) bool { return code >= 200 && code <= 299 }

func IsRedirect(code int) bool { return code >= 300 && code <= 399 }

func IsClientError(code int) bool { return code >= 400 && code <= 499 }

func IsServerError(code int) bool { return code >= 500 && code <= 599 }
