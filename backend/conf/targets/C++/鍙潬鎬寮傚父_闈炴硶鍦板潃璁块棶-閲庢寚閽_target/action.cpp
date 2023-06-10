#include <string>

using namespace std;

char *sqlParamPtr = NULL;
size_t sqlParamSize = 0;
 
void bind(std::size_t pos, const string& val)
{
    std::string str = val;
	sqlParamSize = str.size();
	sqlParamPtr = const_cast<char*>(str.c_str());
	// do something
}

void execute()
{ 
	if (sqlParamPtr)
	{
	    string param = sqlParamPtr;
	    // do something
	}
}

void bind_and_execute()
{
    string param("abc");
    size_t pos = 0;
    bind(pos,param);
    execute();
}